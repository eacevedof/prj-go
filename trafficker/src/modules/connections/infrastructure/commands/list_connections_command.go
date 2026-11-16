package commands

import (
    "errors"
    "fmt"
    "strconv"
    "strings"

    "trafficker/src/core/config"
    "trafficker/src/modules/connections/application/listconnections"
    "trafficker/src/modules/connections/domain/entities"
    "trafficker/src/modules/connections/domain/enums"
    "trafficker/src/modules/connections/domain/exceptions"
    "trafficker/src/modules/ipowners/application/resolveipowners"
    sharedcommands "trafficker/src/modules/shared/infrastructure/commands"
    "trafficker/src/modules/shared/infrastructure/repositories/configuration"
)

const FORMAT_JSON = "json"

// ListConnectionsCommand ./console connections-list [--filter="established,remote"] [--format=json] [--whois=0|1]
// la foto de ahora por consola: conexiones y escuchas con su proceso y usuario. No envia nada al siem.
//
// Whois (pais y organizacion de la IP remota): por defecto solo en modo local (sin .env, como el localhost-traffic
// de Rust); con .env lo resuelve el backend del siem. --whois lo fuerza en un sentido u otro. Es el controller de los
// dos casos de uso (connections y ipowners): los junta para pintarlos.
type ListConnectionsCommand struct {
    sharedcommands.AbstractConsoleCommand
}

func ListConnectionsCommandGetInstance(input []string) *ListConnectionsCommand {
    return &ListConnectionsCommand{AbstractConsoleCommand: sharedcommands.NewAbstractConsoleCommand(input)}
}

func (c *ListConnectionsCommand) Run() {
    listConnectionsResultDto, err := listconnections.ListConnectionsServiceGetInstance().Invoke(
        listconnections.ListConnectionsDtoFromPrimitives(map[string]any{"filter": c.GetOptionValue("filter")}),
    )
    if err != nil {
        var connectionsException *exceptions.ConnectionsException
        if errors.As(err, &connectionsException) {
            c.ExitWithDomainError(connectionsException)
        }
        c.ExitWithUnexpectedError(err, "ListConnectionsCommand.Run")
    }

    var resolveIpOwnersResultDto *resolveipowners.ResolveIpOwnersResultDto
    if c.isWhoisEnabled() {
        resolveIpOwnersResultDto = c.getResolveIpOwnersResultDtoOrNil(listConnectionsResultDto.GetConnectionEntities())
    }

    if strings.ToLower(c.GetOptionValue("format")) == FORMAT_JSON {
        c.printJson(listConnectionsResultDto, resolveIpOwnersResultDto)
        return
    }
    c.printTable(listConnectionsResultDto, resolveIpOwnersResultDto)
}

// isWhoisEnabled --whois=1 / --whois=0 mandan; sin la opcion, solo en modo local.
func (c *ListConnectionsCommand) isWhoisEnabled() bool {
    switch strings.ToLower(c.GetOptionValue("whois")) {
    case "1", "true", "on", "yes":
        return true
    case "0", "false", "off", "no":
        return false
    }
    return configuration.GetInstance().IsLocalMode()
}

// getResolveIpOwnersResultDtoOrNil nil si falla la cache en disco: se avisa y la lista sale igual, sin dueños.
func (c *ListConnectionsCommand) getResolveIpOwnersResultDtoOrNil(connectionEntities []*entities.ConnectionEntity) *resolveipowners.ResolveIpOwnersResultDto {
    remoteIps := []string{}
    for _, connectionEntity := range connectionEntities {
        if !connectionEntity.GetSocketEntity().IsListening() {
            remoteIps = append(remoteIps, connectionEntity.GetSocketEntity().GetRemoteIp())
        }
    }
    resolveIpOwnersResultDto, err := resolveipowners.ResolveIpOwnersServiceGetInstance().Invoke(
        resolveipowners.ResolveIpOwnersDtoFromPrimitives(map[string]any{"ips": remoteIps}),
    )
    if err != nil {
        c.LogException(err, "ListConnectionsCommand.getResolveIpOwnersResultDtoOrNil")
        fmt.Println(c.GetYellow("whois not available: " + err.Error()))
        return nil
    }
    return resolveIpOwnersResultDto
}

func (c *ListConnectionsCommand) printJson(listConnectionsResultDto *listconnections.ListConnectionsResultDto,
    resolveIpOwnersResultDto *resolveipowners.ResolveIpOwnersResultDto) {
    listConnectionsPrimitives := listConnectionsResultDto.ToArray()
    if resolveIpOwnersResultDto != nil {
        for _, connection := range listConnectionsPrimitives.Connections {
            remoteIp, _ := connection["remote_ip"].(string)
            if ipOwnerEntity := resolveIpOwnersResultDto.GetIpOwnerEntityOrNil(remoteIp); ipOwnerEntity != nil {
                for key, value := range ipOwnerEntity.ToPrimitives() {
                    connection[key] = value
                }
            }
        }
    }
    c.PrintResult(listConnectionsPrimitives)
}

// printTable columnas alineadas; el color se pone despues de rellenar (los codigos ANSI no ocupan columnas).
func (c *ListConnectionsCommand) printTable(listConnectionsResultDto *listconnections.ListConnectionsResultDto,
    resolveIpOwnersResultDto *resolveipowners.ResolveIpOwnersResultDto) {
    headers := []string{"PROTO", "DIRECTION", "STATE", "LOCAL", "REMOTE", "PID:PROCESS", "USER"}
    if resolveIpOwnersResultDto != nil {
        headers = append(headers, "COUNTRY", "ORGANIZATION")
    }
    rows := [][]string{}
    for _, connectionEntity := range listConnectionsResultDto.GetConnectionEntities() {
        row := c.getRow(connectionEntity)
        if resolveIpOwnersResultDto != nil {
            row = append(row, c.getOwnerCells(connectionEntity, resolveIpOwnersResultDto)...)
        }
        rows = append(rows, row)
    }
    widths := make([]int, len(headers))
    for column, header := range headers {
        widths[column] = len(header)
        for _, row := range rows {
            widths[column] = max(widths[column], len([]rune(row[column])))
        }
    }

    fmt.Println(c.GetBlue(c.getPaddedLine(headers, widths)))
    for _, row := range rows {
        fmt.Println(c.getColoredLine(c.getPaddedLine(row, widths), row[1]))
    }

    fmt.Println()
    fmt.Printf("%d sockets · %d listening · %d inbound · %d outbound\n",
        len(rows),
        listConnectionsResultDto.GetCountByDirection(enums.DirectionEnum.LISTEN),
        listConnectionsResultDto.GetCountByDirection(enums.DirectionEnum.INBOUND),
        listConnectionsResultDto.GetCountByDirection(enums.DirectionEnum.OUTBOUND),
    )
    if resolveIpOwnersResultDto != nil {
        c.printWhoisSummary(resolveIpOwnersResultDto)
    }
    if unreadableProcesses := listConnectionsResultDto.GetUnreadableProcesses(); unreadableProcesses > 0 {
        fmt.Println(c.GetYellow(config.Tr(enums.ConnectionsMessageEnum.PROCESSES_PARTIALLY_READABLE, unreadableProcesses)))
    }
    if configuration.GetInstance().IsLocalMode() {
        fmt.Println(c.GetBlue("local mode (no .env): nothing is sent to the siem"))
    }
}

func (c *ListConnectionsCommand) printWhoisSummary(resolveIpOwnersResultDto *resolveipowners.ResolveIpOwnersResultDto) {
    summary := fmt.Sprintf("whois: %d from cache · %d looked up", resolveIpOwnersResultDto.GetFromCache(), resolveIpOwnersResultDto.GetLookedUp())
    if failed := resolveIpOwnersResultDto.GetFailed(); failed > 0 {
        summary += fmt.Sprintf(" · %d failed (retried in 1 h)", failed)
    }
    if pending := resolveIpOwnersResultDto.GetPending(); pending > 0 {
        summary += fmt.Sprintf(" · %d pending (next run)", pending)
    }
    fmt.Println(summary)
}

func (c *ListConnectionsCommand) getRow(connectionEntity *entities.ConnectionEntity) []string {
    socketEntity := connectionEntity.GetSocketEntity()
    protocol := socketEntity.GetProtocol()
    if socketEntity.GetIpVersion() == 6 {
        protocol += "6"
    }
    processLabel := connectionEntity.GetProcessEntity().GetName()
    if !connectionEntity.GetProcessEntity().IsUnknown() {
        processLabel = strconv.Itoa(connectionEntity.GetProcessEntity().GetPid()) + ":" + processLabel
    }
    return []string{
        protocol,
        connectionEntity.GetDirection(),
        socketEntity.GetState(),
        socketEntity.GetLocalEndpoint(),
        socketEntity.GetRemoteEndpoint(),
        processLabel,
        connectionEntity.GetUserName(),
    }
}

// getOwnerCells pais y organizacion; vacio en las escuchas, "-" si no se ha podido saber en esta pasada.
func (c *ListConnectionsCommand) getOwnerCells(connectionEntity *entities.ConnectionEntity,
    resolveIpOwnersResultDto *resolveipowners.ResolveIpOwnersResultDto) []string {
    if connectionEntity.GetSocketEntity().IsListening() {
        return []string{"", ""}
    }
    ipOwnerEntity := resolveIpOwnersResultDto.GetIpOwnerEntityOrNil(connectionEntity.GetSocketEntity().GetRemoteIp())
    if ipOwnerEntity == nil {
        return []string{"-", "-"}
    }
    country := ipOwnerEntity.GetCountry()
    if country == "" {
        country = "-"
    }
    return []string{country, ipOwnerEntity.GetLabel()}
}

func (c *ListConnectionsCommand) getPaddedLine(cells []string, widths []int) string {
    paddedCells := make([]string, len(cells))
    for column, cell := range cells {
        // relleno por runas: una organizacion con tildes no descuadra la tabla
        paddedCells[column] = cell + strings.Repeat(" ", max(0, widths[column]-len([]rune(cell))))
    }
    return strings.TrimRight(strings.Join(paddedCells, "  "), " ")
}

// getColoredLine escuchas en amarillo, entrantes en verde, salientes sin color.
func (c *ListConnectionsCommand) getColoredLine(line string, direction string) string {
    switch direction {
    case enums.DirectionEnum.LISTEN:
        return c.GetYellow(line)
    case enums.DirectionEnum.INBOUND:
        return c.GetGreen(line)
    }
    return line
}
