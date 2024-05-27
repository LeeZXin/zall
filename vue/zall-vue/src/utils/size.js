const b = 1;
const kb = 1024 * b;
const mb = 1024 * kb;
const gb = 1024 * mb;
const tb = 1024 * gb;

const KB = "KB";
const MB = "MB";
const GB = "GB";
const TB = "TB";

const readableVolumeSize = size => {
    let fixed = 1;
    let ret = calcUnit(size);
    if (ret.size === 0) {
        fixed = 0
    }
    return ret.size.toFixed(fixed) + ret.unit.unit;
}

const calcUnit = size => {
    let base;
    let unit;
    if (size >= tb) {
        base = tb;
        unit = new Unit(TB);
    } else if (size >= gb) {
        base = gb;
        unit = new Unit(GB);
    } else if (size >= mb) {
        base = mb;
        unit = new Unit(MB);
    } else {
        base = kb;
        unit = new Unit(KB);
    }
    return {
        size: size / base,
        unit
    };
}

class Unit {
    unit;
    toNumber(size) {
        switch (this.unit) {
            case KB:
                return size * kb
            case MB:
                return size * mb
            case GB:
                return size * gb
            case TB:
                return size * tb
            default:
                return size
        }
    }
    constructor(unit) {
        this.unit = unit;
    }
}

export {
    readableVolumeSize,
    calcUnit,
    Unit
}