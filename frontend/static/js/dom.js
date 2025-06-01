const allowedMimeTypes = [
    "application/pdf",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain"
];

export const validFileType = (file) => {
    return allowedMimeTypes.includes(file.type);
}

export const displayFileInfo = (file) => {
    // Format file size
    const size = formatFileSize(file.size);

    // Get file icon based on type
    const iconType = getFileIconType(file);

    filePreview.innerHTML = `
      <div class="file-metadata">
        <div class="file-icon ${iconType}"></div>
        <div class="file-info">
          <p><strong>File name:</strong> ${file.name}</p>
          <p><strong>File size:</strong> ${size}</p>
          <p><strong>Last modified:</strong> ${new Date(
        file.lastModified
    ).toLocaleDateString()}</p>
        </div>
      </div>
    `;
}

const formatFileSize = (bytes) => {
    const sizes = ["Bytes", "KB", "MB", "GB"];
    if (bytes === 0) return "0 Byte";
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return (bytes / Math.pow(1024, i)).toFixed(1) + " " + sizes[i];
}

const getFileIconType = (file) => {
    const type = file.type || "";
    if (type.includes("pdf")) return "pdf-icon";
    if (type.includes("word")) return "docx-icon";
    if (type.includes("text")) return "txt-icon";
    return "default-icon";
}
