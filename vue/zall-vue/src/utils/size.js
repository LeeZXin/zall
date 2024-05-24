const b = 1;
const kb = 1024 * b;
const mb = 1024 * kb;
const gb = 1024 * mb;
const tb = 1024 * gb;

const readableVolumeSize = size => {
    let base;
    let unit;
    if (size >= tb) {
        base = tb;
        unit = "TB";
    } else if (size >= gb) {
        base = gb;
        unit = "GB";
    } else if (size >= mb) {
        base = mb;
        unit = "MB";
    } else if (size >= kb) {
        base = kb;
        unit = "KB";
    } else {
        base = b;
        unit = "B";
    }
    let ret = size / base;
    if (ret <= 0) {
        return "";
    }
    return ret.toFixed(1) + unit;
}

export {
    readableVolumeSize
}