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
    const usernameValue = username.value;
    const emailValue = email.value.trim();
    const passwordValue = password.value;

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
    const regex = /^(?!.*\s)[a-zA-Z\d_]{3,20}$/

    let isValid = true

    if (usernameValue === '') {
        setError(username, "Username cannot be empty")
        isValid = false

    } else if (usernameValue.length < 3) {
        setError(username, "Username must contain at least 3 characters")
        isValid = false

    } else if (!regex.test(usernameValue)) {
        setError(username, "Username must only contain alphabets, numbers, and/or _")
        isValid = false

    }
    else {
        setSuccess(username)
    }

    return isValid
}

function validateEmail(emailValue) {
    const regex = /^(?!.*\s)[\w-\.]+@([\w-]+\.)+[\w-]{2,}$/

    let isValid = true

    if (!regex.test(emailValue)) {
        setError(email, "Invalid email")
        isValid = false

    } else {
        setSuccess(email)
    }

    return isValid
}

function validatePassword(passwordValue) {
    const regex = /(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,128}/

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
    
    errorDisplay.animate([
        { opacity: 0 },
        { opacity: 1 }
    ], {
        duration: 300,
        iterations: 1
    })

    containerInput.classList.add('error');
    errorDisplay.innerText = message;
}

const setSuccess = (element) => {
    const containerInput = element.parentElement;
    const errorDisplay = containerInput.querySelector('.error');

    containerInput.classList.remove('error');
    errorDisplay.innerText = '';
}

// TODO
// confirm password
// name should not contain more than one space in bwteen names
