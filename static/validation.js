const form = document.getElementById('form')
const fullName = document.getElementById('name')
const username = document.getElementById('username')
const email = document.getElementById('email')
const password = document.getElementById('password')

form.addEventListener('submit', e => {
    e.preventDefault();
    const isFormValid = validateInputs();

    if (isFormValid) {
        form.submit()
    }

})

function validateInputs() {
    let isValid = true

    const fullNameValue = fullName.value.trim();
    const usernameValue = username.value.trim();
    const emailValue = email.value.trim();
    const passwordValue = password.value.trim();

    if (fullNameValue === '') {
        setError(fullName, "Please enter your name")
        isValid = false

    } else if (!isValidName(fullNameValue)) {
        setError(fullName, "Please enter a valid full name")
        isValid = false

    } else {
        setSuccess(fullName)
    }

    if (usernameValue === '') {
        setError(username, "Username cannot be empty")
        isValid = false

    } else if (usernameValue.length < 3) {
        setError(username, "Username must contain at least 3 characters")
        isValid = false

    } else if (!isValidUsername(usernameValue)) {
        setError(username, "Username must only caontain alphabets, numbers, and/or _")
        isValid = false

    } else {
        setSuccess(username)
    }

    if (!isValidEmail(emailValue)) {
        setError(email, "Inalid email")
        isValid = false

    } else {
        setSuccess(email)
    }

    if (passwordValue.length < 8) {
        setError(password, "Password must be at least 8 characters long")
        isValid = false

    } else if (!isValidPassword(passwordValue)) {
        setError(password, "Password must caontain lower-case letter, upper-case letter, and number")
        isValid = false

    } else {
        setSuccess(password)
    }

    return isValid
}

const setError = (element, message) => {
    const containerInput = element.parentElement;
    const errorDisplay = containerInput.querySelector('.error');

    containerInput.classList.add('error');
    // containerInput.classList.remove('success');
    errorDisplay.innerText = message;
}

const setSuccess = (element) => {
    const containerInput = element.parentElement;
    const errorDisplay = containerInput.querySelector('.error');

    // containerInput.classList.add('success');
    containerInput.classList.remove('error');
    errorDisplay.innerText = '';
}

function isValidName(nameValue) {
    const regex = /^[a-zA-Z ]{3,100}$/
    return regex.test(nameValue)
}

function isValidUsername(usernameValue) {
    const regex = /^[a-zA-Z\d_]{3,15}$/
    return regex.test(usernameValue)
}

function isValidEmail(emailValue) {
    const regex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,}$/
    return regex.test(emailValue)
}

function isValidPassword(passwordValue) {
    const regex = /(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,64}/
    return regex.test(passwordValue)
}
