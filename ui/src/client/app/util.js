
export const rand = () => Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);

export const guid =() => {
    return rand() + rand() + '-' + rand() + '-' + rand() + '-' +
        rand() + '-' + rand() + rand() + rand();
};

export const isValidURL = (str) => {
    const a  = document.createElement('a');
    a.href = str;
    return a.host;
};