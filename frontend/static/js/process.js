const process = async (form) => {
    const response = await fetch("/upload", {
        method: "POST",
        body: form,
        headers: {
            "Accept": "application/json"
        }
    });

    const data = await response.json();

    return data;
};
