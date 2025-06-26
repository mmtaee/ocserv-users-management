interface Authorization {
    authorization: string;
}

function getAuthorization(): Authorization {
    const token = localStorage.getItem("token");
    return {
        authorization: `Bearer ${token ?? ""}`,
    };
}

export { getAuthorization };