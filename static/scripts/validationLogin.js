const formLogin = document.getElementById('login-form')
const usernameLogin = document.getElementById('login-username')
const passwordLogin = document.getElementById('login-password')

formLogin.addEventListener('submit', e => {
    e.preventDefault();
    const isFormValid = validateInputs();
    
    if (isFormValid) {
        formLogin.submit()
    }
    
})

function validateInputs() {
    
    let isValidForm = true

    const isValidUsername = validateUsername()
    const isValidPassword = validatePassword()

    if (!isValidUsername || !isValidPassword) {
        isValidForm = false
    }
    
    return isValidForm
}

function validateUsername() {
    const usernameValue = usernameLogin.value.trim()

    let isValid = true

    if (usernameValue.length < 3 || usernameValue.length > 15) {
        setError()
        isValid = false
    }

    return isValid
}

function validatePassword() {
    const passwordValue = passwordLogin.value.trim()

    let isValid = true

    if (passwordValue.length < 8 || passwordValue.length > 64) {
        setError()
        isValid = false
    }

    return isValid
}

function setError() {
    const errorDisplay = document.querySelector('.login-error')
    errorDisplay.classList.remove('hidden')
    errorDisplay.innerText = "Invalid username/password"
}