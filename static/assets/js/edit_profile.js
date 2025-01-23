function setEditProfileError(msg) {
	htmx.find("#edit-profile-error").innerHTML = msg
}

function clearEditProfileError() {
	htmx.find("#edit-profile-error").innerHTML = ""
}

function showNewProfileImage(elt, e) {
	const file = elt.files[0]
	console.log(file)
	console.log(file.type)
	if (file.type.indexOf("image/") !== 0) {
		setEditProfileError("profile photo must be an image")
	} else if (file.size >= 1000000) {
		setEditProfileError("profile photo must be less than 1mb")
	} else {
		clearEditProfileError() 
	} 
	const fileReader = new FileReader()
	fileReader.readAsDataURL(file)
	const imageElt = htmx.find("#profile-photo-container img")
	htmx.on(fileReader, "load", () => {
		imageElt.src = fileReader.result
		imageElt.className = "profile-photo"
	})
}
