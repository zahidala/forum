document.addEventListener("DOMContentLoaded", function () {
    // Get the dropdown elements
  const dropdownToggle = document.querySelector(".dropdown-toggle");
  const dropdownMenu = document.querySelector(".dropdown-menu");
  const checkboxes = dropdownMenu.querySelectorAll('input[type="checkbox"]');

  let selectedOptions = [];

  // Toggle dropdown menu on button click
  dropdownToggle.addEventListener("click", function (event) {
    event.stopPropagation(); // Prevent click from propagating to document
    dropdownMenu.classList.toggle("open");
    dropdownMenu.style.display = "block";

    if (dropdownToggle.classList.contains("dropdown-error")) {
      dropdownToggle.classList.remove("dropdown-error");
    }

    // close if already open

    if (dropdownMenu.classList.contains("open")) {
      dropdownMenu.style.display = "none";
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

//   if (!selectedOptions.length) {
//     dropdownToggle.classList.add("dropdown-error");
//     return;
//   }
})