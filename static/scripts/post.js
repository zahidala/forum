document.addEventListener("DOMContentLoaded", () => {
  // show uploaded images in post from the server

  const postId = document.getElementById("post-id").value;
  const postAttachmentsContainer = document.getElementById(`post-images-${postId}`);
  const postAttachments = document.getElementById(`post-attachments-${postId}`).value.split(",").filter((attachment) => attachment !== "");

  postAttachments.forEach((attachment) => {
    const img = document.createElement("img");
    img.src = attachment;
    img.alt = "Post Attachment";
    img.width = 100;
    const fileNameContainer = document.createElement("div");
    fileNameContainer.className = "file-name-container";
    const fileNameSpan = document.createElement("span");
    fileNameSpan.textContent = attachment.split("/").pop();

    fileNameContainer.appendChild(fileNameSpan);
    const uploadedAttachmentDiv = document.createElement("div");
    uploadedAttachmentDiv.className = "uploaded-image";
    uploadedAttachmentDiv.appendChild(img);
    uploadedAttachmentDiv.appendChild(fileNameContainer);

    postAttachmentsContainer.appendChild(uploadedAttachmentDiv);

    postAttachmentsContainer.style.display = "flex";
  });
});