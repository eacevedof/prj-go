package resolveipowners

import (
    "net/netip"
    "sync"
    "time"

    "trafficker/src/modules/ipowners/domain/entities"
    "trafficker/src/modules/ipowners/domain/enums"
    "trafficker/src/modules/ipowners/domain/services"
    "trafficker/src/modules/ipowners/infrastructure/repositories"
    "trafficker/src/modules/shared/infrastructure/repositories/configuration"
)

// ResolveIpOwnersService de quien es cada IP remota (whois por RDAP), con cache en disco:
//   - IP no publicas (privada, loopback, link-local, CGNAT...): se nombran en local, sin preguntar.
//   - En la cache, un rango vigente que la contenga: de ahi (una consulta sirve para todo el bloque).
//   - El resto se pregunta, como mucho MAX_LOOKUPS_PER_RUN por pasada y PARALLEL_LOOKUPS a la vez; un fallo se
//     recuerda una hora. Lo nuevo se guarda en la cache, sin lo caducado.
//
// Un fallo de red o de RDAP nunca es un error del servicio: esa IP sale sin dueño. Solo el disco (cache) da error.
type ResolveIpOwnersService struct {
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
    ipOwnersCacheJsonRepository    *repositories.IpOwnersCacheJsonRepository
    ipOwnersReaderRdapRepository   *repositories.IpOwnersReaderRdapRepository
    ipScopeResolver                *services.IpScopeResolver
    ipOwnerFinder                  *services.IpOwnerFinder
    rdapResponseParser             *services.RdapResponseParser
}

func ResolveIpOwnersServiceGetInstance() *ResolveIpOwnersService {
    environmentReaderRawRepository := configuration.GetInstance()
    return &ResolveIpOwnersService{
        environmentReaderRawRepository: environmentReaderRawRepository,
        ipOwnersCacheJsonRepository:    repositories.IpOwnersCacheJsonRepositoryGetInstance(),
        ipOwnersReaderRdapRepository:   repositories.IpOwnersReaderRdapRepositoryGetInstance(environmentReaderRawRepository.GetAppVersion()),
        ipScopeResolver:                services.IpScopeResolverGetInstance(),
        ipOwnerFinder:                  services.IpOwnerFinderGetInstance(),
        rdapResponseParser:             services.RdapResponseParserGetInstance(),
    }
}

func (s *ResolveIpOwnersService) Invoke(resolveIpOwnersDto *ResolveIpOwnersDto) (*ResolveIpOwnersResultDto, error) {
    statePath := s.environmentReaderRawRepository.GetStatePath()
    cachedIpOwnerEntities, err := s.ipOwnersCacheJsonRepository.GetIpOwnerEntities(statePath)
    if err != nil {
        return nil, err
    }

    now := time.Now().UTC()
    resultDto := &ResolveIpOwnersResultDto{ipOwnerEntityByIp: map[string]*entities.IpOwnerEntity{}}
    addrsToLookUp := []netip.Addr{}
    for _, addr := range resolveIpOwnersDto.GetAddrs() {
        if scope := s.ipScopeResolver.GetScope(addr); scope != enums.IpScopeEnum.PUBLIC {
            resultDto.ipOwnerEntityByIp[addr.String()] = entities.LocalIpOwnerEntity(addr, scope)
            continue
        }
        if cachedIpOwnerEntity := s.ipOwnerFinder.GetIpOwnerEntityOrNil(cachedIpOwnerEntities, addr, now); cachedIpOwnerEntity != nil {
            resultDto.ipOwnerEntityByIp[addr.String()] = cachedIpOwnerEntity
            resultDto.fromCache++
            continue
        }
        addrsToLookUp = append(addrsToLookUp, addr)
    }
    if len(addrsToLookUp) > enums.IpOwnersRdapEnum.MAX_LOOKUPS_PER_RUN {
        resultDto.pending = len(addrsToLookUp) - enums.IpOwnersRdapEnum.MAX_LOOKUPS_PER_RUN
        addrsToLookUp = addrsToLookUp[:enums.IpOwnersRdapEnum.MAX_LOOKUPS_PER_RUN]
    }
    if len(addrsToLookUp) == 0 {
        return resultDto, nil
    }

    newIpOwnerEntities := s.lookUpAll(addrsToLookUp, now, resultDto)
    validIpOwnerEntities := s.ipOwnerFinder.GetValidIpOwnerEntities(append(cachedIpOwnerEntities, newIpOwnerEntities...), now)
    if err := s.ipOwnersCacheJsonRepository.SaveIpOwnerEntities(statePath, validIpOwnerEntities); err != nil {
        return nil, err
    }
    return resultDto, nil
}

// lookUpAll pregunta las IP de PARALLEL_LOOKUPS en PARALLEL_LOOKUPS. Antes de preguntar cada una mira si una
// respuesta de esta misma pasada ya cubre su rango (dos IP del mismo bloque = una consulta).
func (s *ResolveIpOwnersService) lookUpAll(addrs []netip.Addr, now time.Time, resultDto *ResolveIpOwnersResultDto) []*entities.IpOwnerEntity {
    var mutex sync.Mutex
    newIpOwnerEntities := []*entities.IpOwnerEntity{}
    addrsChannel := make(chan netip.Addr)
    var waitGroup sync.WaitGroup
    for range enums.IpOwnersRdapEnum.PARALLEL_LOOKUPS {
        waitGroup.Add(1)
        go func() {
            defer waitGroup.Done()
            for addr := range addrsChannel {
                mutex.Lock()
                knownIpOwnerEntity := s.ipOwnerFinder.GetIpOwnerEntityOrNil(newIpOwnerEntities, addr, now)
                if knownIpOwnerEntity != nil && !knownIpOwnerEntity.IsFailed() {
                    resultDto.ipOwnerEntityByIp[addr.String()] = knownIpOwnerEntity
                    resultDto.fromCache++
                    mutex.Unlock()
                    continue
                }
                mutex.Unlock()

                ipOwnerEntity := s.lookUp(addr, now)
                mutex.Lock()
                newIpOwnerEntities = append(newIpOwnerEntities, ipOwnerEntity)
                resultDto.ipOwnerEntityByIp[addr.String()] = ipOwnerEntity
                resultDto.lookedUp++
                if ipOwnerEntity.IsFailed() {
                    resultDto.failed++
                }
                mutex.Unlock()
            }
        }()
    }
    for _, addr := range addrs {
        addrsChannel <- addr
    }
    close(addrsChannel)
    waitGroup.Wait()
    return newIpOwnerEntities
}

func (s *ResolveIpOwnersService) lookUp(addr netip.Addr, now time.Time) *entities.IpOwnerEntity {
    rdapResponse, err := s.ipOwnersReaderRdapRepository.GetRdapResponse(addr)
    if err != nil {
        return entities.FailedIpOwnerEntity(addr, now)
    }
    ipOwnerEntity := s.rdapResponseParser.GetIpOwnerEntityOrNil(rdapResponse, addr, now)
    if ipOwnerEntity == nil {
        return entities.FailedIpOwnerEntity(addr, now)
    }
    return ipOwnerEntity
}
