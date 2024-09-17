document.addEventListener("DOMContentLoaded", function () {
  // wysiwyg editor related code

  function formatText(command) {
    document.execCommand(command, false, null);
  }

  // Get the dropdown elements
  const dropdownMenu = document.querySelector(".dropdown-menu");
  const checkboxes = dropdownMenu.querySelectorAll('input[type="checkbox"]');
  const dropdownSelect = document.querySelector(".dropdown-select");
  const dropdownToggle = dropdownSelect.querySelector(".dropdown-toggle");

  let selectedOptions = [];

  // Toggle dropdown menu on button click
  dropdownToggle.addEventListener("click", function (event) {
    event.stopPropagation(); // Prevent click from propagating to document
    dropdownMenu.classList.toggle("open");

    // close if open

    if (dropdownMenu.classList.contains("open")) {
      dropdownMenu.style.display = "block";
    } else {
      dropdownMenu.style.display = "none";
    }

    if (dropdownToggle.classList.contains("dropdown-error")) {
      dropdownToggle.classList.remove("dropdown-error");
    }
  });

  // Close the dropdown when clicking outside
  document.addEventListener("click", function (event) {
    if (
      !dropdownMenu.contains(event.target) &&
      !dropdownToggle.contains(event.target)
    ) {
      dropdownMenu.classList.remove("open");
      dropdownMenu.style.display = "none";
    }
  });

  // Update button text with selected values
  checkboxes.forEach((checkbox) => {
    checkbox.addEventListener("change", function () {
      updateDropdownText();
    });
  });

  function updateDropdownText() {
    const selectedOptionsLabel = Array.from(checkboxes)
      .filter((checkbox) => checkbox.checked)
      .map((checkbox) => checkbox.parentElement.textContent.trim());

    selectedOptions = Array.from(checkboxes)
      .filter((checkbox) => checkbox.checked)
      .map((checkbox) => +checkbox.value);
  
    if (selectedOptionsLabel.length > 0) {
      dropdownToggle.textContent =
        selectedOptionsLabel.length > 2
          ? `${selectedOptionsLabel.slice(0, 2).join(", ")} +${
              selectedOptionsLabel.length - 2
            } more`
          : selectedOptionsLabel.join(", ");
    } else {
      dropdownToggle.innerHTML = `Select Categories <i
          class="fa-solid fa-chevron-down"
          style="color: white; margin-left: 5px"
        ></i>`;
    }
  }

  const editorButtons = document.querySelectorAll(".editor-button");

  editorButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const command = button.id;
      formatText(command);
      button.classList.toggle("editor-button-active");
    });
  });

  const editor = document.querySelector(".editor-textarea");

  const postTitle = document.getElementById("post-title-input");

  postTitle.addEventListener("input", () => {
    if (postTitle.classList.contains("post-title-error")) {
      postTitle.classList.remove("post-title-error");
    }
  });

  const editorContainer = document.querySelector(".editor");

  let content = "";

  const imageUrls = [];

  if (editor) {
    editor.addEventListener("focus", () => {
      editorContainer.style.border = "1px solid #2daae9";
    });

    editor.addEventListener("blur", () => {
      editorContainer.style.border = "1px solid #474c54";
    });

    document.addEventListener("DOMContentLoaded", function () {
      const content = document.getElementById("content");
      const placeholder = document.getElementById("placeholder");

      content.addEventListener("input", function () {
        if (content.textContent.trim() !== "") {
          placeholder.style.display = "none";
        } else {
          placeholder.style.display = "block";
        }
      });

      content.addEventListener("focus", function () {
        if (content.textContent.trim() === "") {
          placeholder.style.display = "none";
        }
      });

      content.addEventListener("blur", function () {
        if (content.textContent.trim() === "") {
          placeholder.style.display = "block";
        }
      });
    });

    const editorContent = document.getElementById("content");

    // get content of the editor

    editorContent.addEventListener("input", () => {
      content = editorContent.innerHTML;
    });

    const imageButton = document.getElementById("image");
    const imageInput = document.getElementById("image-input");
    const uploadedImagesContainer = document.getElementById("uploaded-images");

    imageButton.addEventListener("click", () => {
      imageInput.click();
    });

    imageInput.addEventListener("change", async () => {
      const file = imageInput.files[0];

      if (!file) return;

      const formData = new FormData();
      formData.append("image", file);

      try {
        const response = await fetch("/upload", {
          method: "POST",
          body: formData,
        });

        if (response.ok) {
          const data = await response.json();
          const imageUrl = data.image.url;
          imageUrls.push(imageUrl);

          // Create a new div for the uploaded image
          const uploadedImageDiv = document.createElement("div");
          uploadedImageDiv.className = "uploaded-image";

          // Create an img element
          const img = document.createElement("img");
          img.src = imageUrl;
          img.alt = "Uploaded Image";
          img.width = 100;

          // Create a div for the file name
          const fileNameContainer = document.createElement("div");
          fileNameContainer.className = "file-name-container";
          const fileNameSpan = document.createElement("span");
          fileNameSpan.textContent = data.image.name;

          // Append the img and file name div to the uploaded image div
          fileNameContainer.appendChild(fileNameSpan);
          uploadedImageDiv.appendChild(img);
          uploadedImageDiv.appendChild(fileNameContainer);

          // Append the uploaded image div to the uploaded images container
          uploadedImagesContainer.appendChild(uploadedImageDiv);

          uploadedImagesContainer.style.display = "flex";
        }
      } catch (error) {
        console.error(error);
      }
    });
  }

  const postReply = document.getElementById("post-reply");

  if (postReply) {
    // send new post to the server
    postReply.addEventListener("click", async () => {
      const userId = document.getElementById("userId").value;
      const images = imageUrls.join(",");

      if (!selectedOptions.length && postTitle.value.trim() === "" && content.textContent.trim() === "") {
        dropdownToggle.classList.add("dropdown-error");
        postTitle.classList.add("post-title-error");
        editorContainer.classList.add("editor-error");
        return;
      }

      if (!selectedOptions.length) { 
        dropdownToggle.classList.add("dropdown-error")
        return;
      }

      if (postTitle.value.trim() === "") {
        postTitle.classList.add("post-title-error");
        return;
      }

      if (content.textContent.trim() === "") {
        editorContainer.classList.add("editor-error");
        return;
      }

      try {
        const response = await fetch(`/new-post`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            content,
            userId,
            title: postTitle.value,
            selectedCategories: selectedOptions,
            images: images ? images : undefined,
          }),
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