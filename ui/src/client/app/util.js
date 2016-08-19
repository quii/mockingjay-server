
export const rand = () => Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);

export const guid = () => `${rand()}-${rand()}-${rand()}-${rand()}`;

export const isValidURL = (str) => {
  const a = document.createElement('a');
  a.href = str;
  return a.host;
};

export function isJSON(value) {
  try {
    JSON.parse(value);
    return true;
  } catch (e) {
    return false;
  }
}
