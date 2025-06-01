export const process = async (form) => {
    const response = await fetch("/upload", {
        method: "POST",
        body: form,
        headers: {
            "Accept": "application/json"
        }
    });

    if (!response.ok) {
        return "";
    }

    let data = await response.json();

    if (data.status == "error") {
        alert(data.message);
        return "";
    }


    return data.units;
};
