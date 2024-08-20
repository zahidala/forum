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
          method: "PUT",
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
          method: "PUT",
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

  const commentLikedByUser = document.querySelectorAll(".comment-liked-by-user");

  if (commentLikedByUser.length > 0) {
    commentLikedByUser.forEach((comment) => {
      if (comment.value === "true") {
        const commentId = comment.dataset.commentId;
        const commentContainer = document.getElementById(`comment-${commentId}`);
        const likeButton = commentContainer.querySelector(".like-comment");

        likeButton.classList.add("post-liked");
        likeButton.classList.remove("like-comment");
        likeButton.classList.add("remove-like-comment");
        likeButton.innerHTML = `<i class="fa-solid fa-thumbs-up" style="margin-right: 5px"></i> Liked`;
      }
    });
  }

  const removeLikeCommentButtons = document.querySelectorAll(".remove-like-comment");

  removeLikeCommentButtons.forEach((removeLikeButton) => {
    removeLikeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;
      const commentId = removeLikeButton.parentNode.id;
      const postId = document.getElementById("postId").value;

      try {
        const response = await fetch(`/comment/${commentId}/remove-like`, {
          method: "PUT",
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

  const commentDislikedByUser = document.querySelectorAll(".comment-disliked-by-user");

  if (commentDislikedByUser.length > 0) {
    commentDislikedByUser.forEach((comment) => {
      if (comment.value === "true") {
        const commentId = comment.dataset.commentId;
        const commentContainer = document.getElementById(`comment-${commentId}`);
        const dislikeButton = commentContainer.querySelector(".dislike-comment");

        dislikeButton.classList.add("post-disliked");
        dislikeButton.classList.remove("dislike-comment");
        dislikeButton.classList.add("remove-dislike-comment");
        dislikeButton.innerHTML = `<i class="fa-solid fa-thumbs-down" style="margin-right: 5px"></i> Disliked`;
      }
    });
  }

  const removeDislikeCommentButtons = document.querySelectorAll(".remove-dislike-comment");

  removeDislikeCommentButtons.forEach((removeDislikeButton) => {
    removeDislikeButton.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;
      const commentId = removeDislikeButton.parentNode.id;
      const postId = document.getElementById("postId").value;

      try {
        const response = await fetch(`/comment/${commentId}/remove-dislike`, {
          method: "PUT",
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