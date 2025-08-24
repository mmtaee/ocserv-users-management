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

const formatDate = (dateString: string | undefined) => {
    let dateTime = formatDateTime(dateString, "")
    return dateTime.split(" ")[0]
}

const formatDateTimeWithRelative = (
    dateString: string | undefined,
    message: string | undefined
): string => {
    if (!dateString) {
        return message || "";
    }

    const formatted = formatDateTime(dateString, message);
    const date = new Date(dateString);
    const now = new Date();

    // Calculate difference in milliseconds
    const diffTime = now.getTime() - date.getTime();

    // Helper to get full year/month/day difference
    const diffYears = now.getFullYear() - date.getFullYear();
    const diffMonths = (now.getFullYear() - date.getFullYear()) * 12 + (now.getMonth() - date.getMonth());
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    let relative = "";

    if (diffDays === 0) {
        relative = "Today";
    } else if (diffDays === 1) {
        relative = "Yesterday";
    } else if (diffDays === -1) {
        relative = "Tomorrow";
    } else if (Math.abs(diffYears) >= 1) {
        if (diffYears > 0) {
            relative = `${diffYears} year${diffYears > 1 ? "s" : ""} ago`;
        } else {
            relative = `in ${Math.abs(diffYears)} year${Math.abs(diffYears) > 1 ? "s" : ""}`;
        }
    } else if (Math.abs(diffMonths) >= 1) {
        if (diffMonths > 0) {
            relative = `${diffMonths} month${diffMonths > 1 ? "s" : ""} ago`;
        } else {
            relative = `in ${Math.abs(diffMonths)} month${Math.abs(diffMonths) > 1 ? "s" : ""}`;
        }
    } else {
        if (diffDays > 1) {
            relative = `${diffDays} days ago`;
        } else if (diffDays < -1) {
            relative = `in ${Math.abs(diffDays)} days`;
        }
    }

    return `${formatted} (${relative})`;
};


export {bytesToGB, formatDateTime, formatDate, formatDateTimeWithRelative}
