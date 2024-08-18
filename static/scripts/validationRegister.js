const registerForm = document.getElementById('register-form')
const fullName = document.getElementById('name')
const username = document.getElementById('username')
const email = document.getElementById('email')
const password = document.getElementById('password')


registerForm.addEventListener('submit', e => {
    e.preventDefault();

    const isFormValid = validateInputs();

    if (isFormValid) {
        registerForm.submit()
    }

})

function validateInputs() {
    const fullNameValue = fullName.value.trim();
    const usernameValue = username.value.trim();
    const emailValue = email.value.trim();
    const passwordValue = password.value.trim();

    let isValidForm = true

    const isValidName = validateName(fullNameValue)
    const isValidUsername = validateUsername(usernameValue)
    const isValidEmail = validateEmail(emailValue)
    const isValidPassword = validatePassword(passwordValue)

    if (!isValidName || !isValidUsername || !isValidEmail || !isValidPassword) {
        isValidForm = false
    }

    return isValidForm
}

function validateName(fullNameValue) {
    const regex = /^[a-zA-Z ]{3,50}$/

    let isValid = true

    if (fullNameValue === '') {
        setError(fullName, "Please enter your name")
        isValid = false

    } else if (!regex.test(fullNameValue)) {
        setError(fullName, "Please enter a valid full name")
        isValid = false

    } else {
        setSuccess(fullName)
    }

    return isValid
}

function validateUsername(usernameValue) {
    const regex = /^[a-zA-Z\d_]{3,15}$/

    let isValid = true

    if (usernameValue === '') {
        setError(username, "Username cannot be empty")
        isValid = false

    } else if (usernameValue.length < 3) {
        setError(username, "Username must contain at least 3 characters")
        isValid = false

    } else if (!regex.test(usernameValue)) {
        setError(username, "Username must only caontain alphabets, numbers, and/or _")
        isValid = false

    }
    else {
        setSuccess(username)
    }

    return isValid
}

function validateEmail(emailValue) {
    const regex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,}$/

    let isValid = true

    if (!regex.test(emailValue)) {
        setError(email, "Inalid email")
        isValid = false

    } else {
        setSuccess(email)
    }

    return isValid
}

function validatePassword(passwordValue) {
    const regex = /(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,64}/

    let isValid = true

    if (passwordValue.length < 8) {
        setError(password, "Password must be at least 8 characters long")
        isValid = false

    } else if (!regex.test(passwordValue)) {
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
    // errorDisplay.style.display = 'block'
    errorDisplay.innerText = message;
}

const setSuccess = (element) => {
    const containerInput = element.parentElement;
    const errorDisplay = containerInput.querySelector('.error');
    
    containerInput.classList.remove('error');
    // errorDisplay.style.display = 'block'
    errorDisplay.innerText = '';
}

// TODO
// confirm password
// name should not contain more than one space in bwteen names
