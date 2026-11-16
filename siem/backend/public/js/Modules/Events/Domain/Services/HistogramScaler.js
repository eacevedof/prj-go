import { EventSeverityEnum } from "/js/Modules/Events/Domain/Enums/EventSeverityEnum.js";

// EventSeverityEnum va de menos a mas grave: el orden de apilado
const SEVERITIES_BOTTOM_UP = Object.freeze(Object.values(EventSeverityEnum));
// el techo del eje Y es uno de estos por una potencia de 10
const NICE_CEILING_STEPS = Object.freeze([1, 2, 5, 10]);
// posiciones de "MM-DD" y "HH:MM" dentro de "YYYY-MM-DD HH:MM:SS"
const MONTH_DAY_START = 5;
const MONTH_DAY_END = 10;
const HOUR_MINUTE_START = 11;
const HOUR_MINUTE_END = 16;
const SECONDS_PER_DAY = 86400;
const X_TICKS_MAX = 7;
const Y_TICKS = 2;
const BAR_GAP_RATIO = 0.18;

/**
 * Geometria del histograma de eventos: de los buckets de la API a barras apiladas por gravedad (de abajo arriba:
 * debug, info, warning, error, critical), todas con la misma escala. El techo del eje Y es un numero "redondo" por
 * encima del maximo, y las etiquetas del eje X (UTC) se reparten sin pisarse. Funcion pura: sin DOM.
 */
export class HistogramScaler {

    /**
     * @returns {HistogramScaler}
     */
    static getInstance() {
        return new HistogramScaler();
    }

    /**
     * @param {{bucket_at: string, total: number, debug: number, info: number, warning: number, error: number, critical: number}[]} buckets
     * @param {number} bucketSeconds
     * @param {{width: number, height: number, left: number, right: number, top: number, bottom: number}} box tamaño del
     *     dibujo y margenes (ejes y etiquetas)
     * @returns {{bars: Object[], xTicks: {x: number, label: string}[], yTicks: {y: number, label: string}[],
     *     yMax: number, total: number, plot: {x: number, y: number, width: number, height: number}}}
     */
    getChart(buckets, bucketSeconds, box) {
        const plot = {
            x: box.left,
            y: box.top,
            width: Math.max(1, box.width - box.left - box.right),
            height: Math.max(1, box.height - box.top - box.bottom),
        };
        const maxTotal = buckets.reduce((max, bucket) => Math.max(max, Number(bucket.total) || 0), 0);
        const yMax = this.getNiceCeiling(maxTotal);
        const slotWidth = buckets.length ? plot.width / buckets.length : plot.width;
        const barWidth = Math.max(1, slotWidth * (1 - BAR_GAP_RATIO));

        const bars = buckets.map((bucket, index) => {
            let stackedCount = 0;
            const segments = [];
            SEVERITIES_BOTTOM_UP.forEach((severity) => {
                const count = Number(bucket[severity]) || 0;
                if (count === 0) return;
                const segmentHeight = count / yMax * plot.height;
                stackedCount += count;
                segments.push({
                    severity,
                    count,
                    y: plot.y + plot.height - stackedCount / yMax * plot.height,
                    height: segmentHeight,
                });
            });
            return {
                bucketAt: bucket.bucket_at,
                total: Number(bucket.total) || 0,
                x: plot.x + index * slotWidth + (slotWidth - barWidth) / 2,
                width: barWidth,
                slotX: plot.x + index * slotWidth,
                slotWidth,
                segments,
            };
        });

        return {
            bars,
            xTicks: this.#getXTicks(buckets, bucketSeconds, plot, slotWidth),
            yTicks: this.#getYTicks(yMax, plot),
            yMax,
            total: buckets.reduce((sum, bucket) => sum + (Number(bucket.total) || 0), 0),
            plot,
        };
    }

    /**
     * @param {number} value
     * @returns {number} 1, 2, 5 o 10 por una potencia de 10: el primero que no queda por debajo (minimo 1)
     */
    getNiceCeiling(value) {
        if (!(value > 0)) return 1;
        const magnitude = 10 ** Math.floor(Math.log10(value));
        const candidate = NICE_CEILING_STEPS.find((step) => step * magnitude >= value);
        return candidate * magnitude;
    }

    #getYTicks(yMax, plot) {
        const yTicks = [];
        for (let step = 0; step <= Y_TICKS; step++) {
            const value = yMax * step / Y_TICKS;
            // con un techo de 1 no hay "0,5 eventos"
            if (!Number.isInteger(value)) continue;
            yTicks.push({ y: plot.y + plot.height - value / yMax * plot.height, label: String(value) });
        }
        return yTicks;
    }

    #getXTicks(buckets, bucketSeconds, plot, slotWidth) {
        if (buckets.length === 0) return [];
        const spanSeconds = buckets.length * (Number(bucketSeconds) || 0);
        const every = Math.max(1, Math.ceil(buckets.length / X_TICKS_MAX));
        const xTicks = [];
        for (let index = 0; index < buckets.length; index += every) {
            xTicks.push({
                x: plot.x + index * slotWidth + slotWidth / 2,
                label: this.#getTimeLabel(buckets[index].bucket_at, Number(bucketSeconds) || 0, spanSeconds),
            });
        }
        return xTicks;
    }

    // "YYYY-MM-DD HH:MM:SS" -> tramos de un dia o mas: "MM-DD"; rango de mas de un dia: "MM-DD HH:MM"; si no, "HH:MM"
    #getTimeLabel(bucketAt, bucketSeconds, spanSeconds) {
        const text = String(bucketAt ?? "");
        const monthDay = text.slice(MONTH_DAY_START, MONTH_DAY_END);
        const hourMinute = text.slice(HOUR_MINUTE_START, HOUR_MINUTE_END);
        if (bucketSeconds >= SECONDS_PER_DAY) return monthDay;
        if (spanSeconds > SECONDS_PER_DAY) return `${monthDay} ${hourMinute}`;
        return hourMinute;
    }
}
