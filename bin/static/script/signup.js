function validate() {
		var pwd = document.forms["signupForm"]["pwd"].value;
		var rpwd = document.forms["signupForm"]["rpwd"].value;
		var mail = document.forms["signupForm"]["mail"].value;

		if (pwd !== rpwd) {
				alert("Passwords do not match. Please fill again");
				return false;
		}

		if (!(/.*@[a-z|A-Z|0-9]+[.com|.org|.edu]/.test(mail))) {
				alert("Mail id invalid. Please follow proper format.");
				return false;
		}
}
