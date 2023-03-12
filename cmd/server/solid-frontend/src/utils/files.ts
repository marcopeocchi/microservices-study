export const isImage = (str: string) => str.match(/(\.jpg)|(\.jpeg)|(\.png)|(\.webp)|(\.avif)|(\.bmp)|(\.gif)/gmi)

export const isOrderedByDate = (signal: string) => signal === "date"
export const isOrderedByName = (signal: string) => signal === "name"