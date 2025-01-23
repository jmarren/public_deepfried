


// retrieves the dropped files, adds them to the fileInput, and invokes HandleFileInput changed
 function HandleDrop(e) {
	htmx.find('#file-input').files = e.dataTransfer.files
	const fileInput = htmx.find('#file-input')
	HandleFileInputChanged(fileInput)
 }



//  Not sure if this is being used?
function addInput(e, elt) {
	      const parent = elt.parentNode
	      const currentTags = htmx.findAll(elt.parentElement, ".tag-input")
	      if (currentTags.length >= 3) {
		      return
	      }
	      const newNode = document.createElement('input')
	      newNode.type = 'text'
	      newNode.placeholder = 'jazz'
	      newNode.className = 'tag-input'
	      elt.parentNode.insertBefore(newNode, elt)
	      htmx.process(newNode)
      }


// sets an error in the error div
function setUploadErr(msg) {
      const errorElt = htmx.find("#upload-error")
      errorElt.innerHTML = msg
}


/* upload_form.js */
function AddTagButton(elt) {
  const container = htmx.closest(elt, ".input-container")
  const numTagInputs = htmx.findAll(container, "input").length
  if (numTagInputs >= 3) return
  const newTagContainer = document.createElement("div")
  newTagContainer.className="tag-input-container"
  const newInput = document.createElement("input")
  newInput.name=`tag-${numTagInputs + 1}`
  newInput.className="tag-input input-field"
  newInput.placeholder = "jazz"
  newInput.type = "text"
  newTagContainer.appendChild(newInput)
  const newRemoveButton = document.createElement("button")
  newRemoveButton.onclick = (e) => e.target.closest(".tag-input-container").remove()
  newRemoveButton.className="remove-tag-button"
  const removeIcon = document.createElement("div")
  removeIcon.className="remove-icon"
  newRemoveButton.appendChild(removeIcon)
  newTagContainer.appendChild(newRemoveButton)
  container.insertBefore(newTagContainer, elt)
  htmx.process(container)
}

// when the input of the input[name=stems] element changes,
// addStemFile is invoked.
// It will get the names of the stem files and add them to the input[name=stem_file_names] input
// so they can be posted when the upload is confirmed
function addStemFile(elt) {
  console.log("files: ", elt.files)
  const files = elt.files
  const stemsSection = htmx.find("#stems-input-section")
  const container = htmx.find(stemsSection, ".input-container")
  const addButton = htmx.find(stemsSection, "#add-stems-button")

  const stemFileNamesInput = htmx.find("input[name=stem_file_names]") 
  const stemFileNames = JSON.parse(stemFileNamesInput.value)
  console.log(stemFileNames)

  for (let i = 0; i < files.length; i++) {
    const fileElt = document.createElement('div');
    fileElt.innerHTML = files[i].name
    stemFileNames.push(files[i].name)
    fileElt.className = "stem-file"
    container.insertBefore(fileElt, addButton)
    InitUpload(files[i], "stems")
  }
  const newValue = JSON.stringify(stemFileNames)
  console.log(newValue)
  stemFileNamesInput.value = newValue
}



