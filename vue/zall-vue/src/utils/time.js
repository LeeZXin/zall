import moment from 'moment';
import i18n from "../language/i8n"
const t = i18n.global.t;

const second = 1;
const minute = 60 * second;
const hour = 60 * minute;
const day = 24 * hour;
const month = 30 * day;
const year = 12 * month;

const readableTimeComparingNow = dateTimeStr => {
    let sub = moment().unix() - moment(dateTimeStr, "YYYY-MM-DD HH:mm:ss").unix();
    if (sub >= year) {
        return parseInt(sub / year) + t("yearBefore");
    }
    if (sub >= month) {
        return parseInt(sub / month) + t("monthBefore");
    }
    if (sub >= day) {
        return parseInt(sub / day) + t("dayBefore");
    }
    if (sub >= hour) {
        return parseInt(sub / hour) + t("hourBefore");
    }
    if (sub >= minute) {
        return parseInt(sub / minute) + t("minuteBefore");
    }
    return sub + t("secondBefore");
}

const yy = () => {
    console.log(day, t("ss"));
}

export {
    readableTimeComparingNow,
    yy
}