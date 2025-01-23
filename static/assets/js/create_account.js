
function setCreateProfileError(msg) {
	htmx.find("#create-profile-error").innerHTML = msg
}

function clearCreateProfileError() {
	htmx.find("#create-profile-error").innerHTML = ""
}

function showSelectedProfileImage(elt, e) {
	const file = elt.files[0]
	// console.log(file)
	// console.log(file.type)
	if (file.type.indexOf("image/") !== 0) {
		setCreateProfileError("profile photo must be an image")
	} else if (file.size >= 1000000) {
		setCreateProfileError("profile photo must be less than 1mb")
	} else {
		clearCreateProfileError() 
	} 
	const fileReader = new FileReader()
	fileReader.readAsDataURL(file)
	const imageElt = htmx.find("#selected-photo")
	htmx.on(fileReader, "load", () => {
		imageElt.src = fileReader.result
		imageElt.className = "profile-photo"
	})
}
