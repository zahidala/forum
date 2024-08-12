document.addEventListener('DOMContentLoaded', function() {
   // wysiwyg editor related code

   function formatText(command) {
    document.execCommand(command, false, null);
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

  const editorContainer = document.querySelector(".editor");

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

  const linkButton = document.getElementById("link");

  linkButton.addEventListener("click", (e) => {
    const linkDialog = document.getElementById("link-dialog");

    // if editor is focused, close the dialog
    editor.addEventListener("focus", () => {
      linkDialog.classList.add("link-dialog-closing");
      linkButton.classList.remove("editor-button-active");
      setTimeout(() => {
        linkDialog.close();
        document.getElementById('url').value = '';
        document.getElementById('text').value = '';
        linkDialog.classList.remove("link-dialog-closing");
      }, 200);
    });

    if (linkDialog.open) {
      // add the closing class to animate the dialog
      linkDialog.classList.add("link-dialog-closing");

      // close the dialog after the animation completes
      setTimeout(() => {
        linkDialog.close();
        linkDialog.classList.remove("link-dialog-closing");
      }, 200);
    } else {
      linkDialog.show();

      // animate the dialog
      linkDialog.animate(
        [
          {
            transform: "scale(0)",
          },
          {
            transform: "scale(1)",
          },
        ],
        {
          duration: 200,
          easing: "ease",
        }
      );
    }
  });

  const insertLink = document.getElementById("insert-link");

  insertLink.addEventListener("click", () => {
    const url = document.getElementById('url').value;
    const text = document.getElementById('text').value;
    const content = document.getElementById('content');

    content.focus();

    // Create a range and select the text
    const selection = window.getSelection();
    if (selection.rangeCount > 0) {
        const range = selection.getRangeAt(0);
        range.deleteContents();

        // Create a new text node with the provided text
        const textNode = document.createTextNode(text);

        // Insert the text node into the range
        range.insertNode(textNode);

        // Select the text node
        range.selectNode(textNode);
        selection.removeAllRanges();
        selection.addRange(range);

        // Create the link
        document.execCommand('createLink', false, url);

        // Add the link class to the link

        const link = document.querySelector('a[href="' + url + '"]');

        if (link) link.classList.add('editor-link');
    }

    // Close the dialog
    document.getElementById('link-dialog').close();

    // Clear the input fields
    document.getElementById('url').value = '';
    document.getElementById('text').value = '';
});

  const editorContent = document.getElementById("content");

  // get content of the editor

  let content = "";

  editorContent.addEventListener("input", () => {
    content = editorContent.innerHTML;
  });

  const postReply = document.getElementById("post-reply");

  postReply.addEventListener("click", async () => {
    const postId = document.getElementById("postId").value;
    const userId = document.getElementById("userId").value;
    try {
      const response = await fetch(`/post/${postId}/comment`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          content,
          userId,
        }),
      });
      
      if (response.redirected) {
        window.location.href = response.url;
      }
    } catch (error) {
      console.error(error);
    }
  });
});