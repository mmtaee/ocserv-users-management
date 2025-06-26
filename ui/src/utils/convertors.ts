const bytesToGB = (bytes: number): string => {
    return (bytes / (1024 ** 3)).toFixed(6); // returns GB as a string with 6 decimal places
}

const formatDateTime = (dateString: string | undefined, message: string | undefined): string => {
    if (!dateString) {
        return message || ""
    }
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

const formatDay = (dateString: string | undefined) => {
    let dateTime = formatDateTime(dateString, "")
    return dateTime.split(" ")[0]
}
export {bytesToGB, formatDateTime, formatDay}
