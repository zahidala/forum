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

  const likeButton = document.getElementById("post-like");
  const dislikeButton = document.getElementById("post-dislike");

  if (likeButton) { 
    likeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;

      try {
        const response = await fetch(`/post/${postId}/like`, {
          method: "PUT",
          body: JSON.stringify({ userId }),
        });

        if (response.redirected) {
          window.location.href = response.url;
        }
      } catch (error) {
        console.error(error);
      }
    });
  }

  if (dislikeButton) {
    dislikeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;

      try {
        const response = await fetch(`/post/${postId}/dislike`, {
          method: "PUT",
          body: JSON.stringify({ userId }),
        });

        if (response.redirected) {
          window.location.href = response.url;
        }
      } catch (error) {
        console.error(error);
      }
    });
  }

  const removeLikeButton = document.getElementById("post-remove-like");

  if (removeLikeButton) {
    removeLikeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;

      try {
        const response = await fetch(`/post/${postId}/remove-like`, {
          method: "PUT",
          body: JSON.stringify({ userId }),
        });

        if (response.redirected) {
          window.location.href = response.url;
        }
      } catch (error) {
        console.error(error);
      }
    });
  }
});