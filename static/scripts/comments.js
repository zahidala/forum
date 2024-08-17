document.addEventListener('DOMContentLoaded', function() {
  // show uploaded images in comments from the server

  const commentsAttachments = document.querySelectorAll(".comment-attachments");
  const commentImages = document.querySelectorAll(".comment-images");

  commentsAttachments.forEach((commentAttachment, index) => {
    const images = commentAttachment.value.split(",").filter((image) => image !== "");
    images.forEach((image) => {
      const img = document.createElement("img");
      img.src = image;
      img.alt = "Comment Image";
      img.width = 100;
      const fileNameContainer = document.createElement("div");
      fileNameContainer.className = "file-name-container";
      const fileNameSpan = document.createElement("span");
      fileNameSpan.textContent = image.split("/").pop();

      fileNameContainer.appendChild(fileNameSpan);
      const uploadedImageDiv = document.createElement("div");
      uploadedImageDiv.className = "uploaded-image";
      uploadedImageDiv.appendChild(img);
      uploadedImageDiv.appendChild(fileNameContainer);

      // Append the image to the corresponding comment's image container
      if (commentImages[index]) {
        commentImages[index].appendChild(uploadedImageDiv);
        commentImages[index].style.display = "flex";
      }
    });
  });

  // Like button logic

  const likeButtons = document.querySelectorAll(".like-comment");

  likeButtons.forEach((likeButton) => {
    likeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;
      const commentId = likeButton.parentNode.id;
      const postId = document.getElementById("postId").value;

      try {
        const response = await fetch(`/comment/${commentId}/like`, {
          method: "POST",
          body: JSON.stringify({ userId, postId }),
        });

        if (response.redirected) {
          window.location.href = response.url;
        }
      } catch (error) {
        console.error(error);
      }
    });
  });

  const dislikeButtons = document.querySelectorAll(".dislike-comment");

  dislikeButtons.forEach((dislikeButton) => {
    dislikeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;
      const commentId = dislikeButton.parentNode.id;
      const postId = document.getElementById("postId").value;

      try {
        const response = await fetch(`/comment/${commentId}/dislike`, {
          method: "POST",
          body: JSON.stringify({ userId, postId }),
        });

        if (response.redirected) {
          window.location.href = response.url;
        }
      } catch (error) {
        console.error(error);
      }
    });
  });
});