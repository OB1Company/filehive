export const ConvertBase64 = (file) => {
    return new Promise((resolve, reject) => {
        const fileReader = new FileReader();
        fileReader.readAsBinaryString(file)
        fileReader.onload = () => {
            resolve(fileReader.result);
        }
        fileReader.onerror = (error) => {
            reject(error);
        }
    })
}

export const ConvertImageToString = async (file) => {
    const blobString = await ConvertBase64(file);
    return btoa(blobString);
}