// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/jmarren/deepfried/util"
)

func UploadForm(name string, uploadId string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<head hx-head=\"append\"><!--\n\t\t<script hx-preserve=\"true\" src={ util.GetStaticSrc(\"js/hx_upload.js\") }></script>\n\t\t<script hx-preserve=\"true\" src={ util.GetStaticSrc(\"js/upload_form.js\") }></script>\n\t\t--><link href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(util.GetStaticSrc("css/upload_form.css"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 14, Col: 55}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" rel=\"stylesheet\" type=\"text/css\"></head><div id=\"selected-files\" hx-on:click=\"event.stopPropagation()\" hx-file-upload-form=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s", name))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 16, Col: 109}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><!--\t\t<input hidden type=\"file\" name=\"stems\" hx-on:input=\"addStemFile(this)\" /> --><div class=\"selected-file\" id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("file_%s", name))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 18, Col: 62}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"selected-file-left\"><div class=\"artwork-input-container\"><div class=\"input-container\"><label for=\"artwork\">select an image</label><!--\n\t\t\t\t\t\t<input\n\t\t\t\t\t\t\thx-file-upload-type=\"artwork\"\n\t\t\t\t\t\t\thx-file-upload-id={ name }\n\t\t\t\t\t\t\ttype=\"file\"\n\t\t\t\t\t\t\tname=\"artwork\"\n\t\t\t\t\t\t\taccept=\"image/*\"\n\t\t\t\t\t\t\thx-on:change=\"HandleArtworkSelected(this)\"\n\t\t\t\t\t\t/>\n\t\t\t\t\t\t--></div><img class=\"file-image\" id=\"artwork-display\" hx-on:click=\"htmx.find(&#39;input[name=artwork]&#39;).click()\"></div><div class=\"file-title\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(name) < 19 {
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 37, Col: 12}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s...", name[0:16]))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 39, Col: 40}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span hx-upload-status></span></div><div class=\"progress-bar\" hx-upload-progress><div id=\"progress-marker\"></div></div></div><div class=\"selected-file-right\"><div class=\"upload-input-form\"><!-- HIDDEN INPUTS --><input hidden type=\"text\" name=\"audio_file_name\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 50, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input hidden type=\"text\" name=\"artwork_file_name\"> <input hidden type=\"text\" name=\"stem_file_names\" value=\"[]\"> <input hidden type=\"text\" name=\"upload-id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(uploadId)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/UploadForm.templ`, Line: 53, Col: 64}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"upload-input-section\"><div class=\"input-container\"><label for=\"title\">title</label> <input class=\"input-field\" type=\"text\" name=\"title\"></div><div class=\"input-container\"><label for=\"bpm\">bpm</label> <input class=\"input-field\" type=\"number\" name=\"bpm\"></div><div class=\"input-container\"><label for=\"musical-key\">key</label> <select class=\"input-field\" name=\"musical-key\"><option value=\"C\">C</option> <option value=\"D\">D</option> <option value=\"E\">E</option> <option value=\"F\">F</option> <option value=\"G\">G</option> <option value=\"A\">A</option> <option value=\"B\">B</option></select></div><div class=\"input-container\"><label for=\"key-signature\">key signature</label> <select class=\"input-field\" name=\"musical-key-signature\"><option value=\"sharp\">#</option> <option value=\"flat\">&#9837;</option> <option value=\"natural\">natural</option></select></div><div class=\"input-container\"><label for=\"major-minor\">Major/Minor</label> <select class=\"input-field\" name=\"major-minor\"><option value=\"Major\">Major</option> <option value=\"Minor\">Minor</option></select></div></div><div class=\"upload-input-section\"><div class=\"input-container\"><label for=\"usage\">usage rights</label> <textarea class=\"input-field\" name=\"usage\" maxlength=\"255\"></textarea></div><div class=\"input-container tags-inputs\"><div class=\"tags-label\">tags</div><div class=\"tag-input-container\"><input class=\"tag-input input-field\" type=\"text\" name=\"tag-1\" placeholder=\"jazz\"></div><button id=\"add-tag-button\" class=\"input-form-button\" hx-on:click=\"AddTagButton(this)\">add tag</button></div></div><div class=\"upload-input-section\" id=\"stems-input-section\"><div class=\"input-container\"><div>stems</div><button id=\"add-stems-button\" class=\"input-form-button\" hx-on:click=\"htmx.find(&#39;input[name=stems-1]&#39;).click()\">add stem file</button></div></div></div><div id=\"upload-error\"></div><button class=\"upload-file-button\" hx-post=\"/audio\" hx-target=\"#upload-error\" hx-swap=\"innerHTML\" hx-include=\"closest .selected-file-right\">post</button></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var uploadFormFile2OnceHandle = templ.NewOnceHandle()

/*

func UploadFormFile2Head() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var9 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var9 == nil {
			templ_7745c5c3_Var9 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var10 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<head><script id=\"upload-form-file-script\">\n\n\t\tfunction AddTagButton(elt) {\n\t\t\tconst container = htmx.closest(elt, \".input-container\")\n\t\t\tconst numTagInputs = htmx.findAll(container, \"input\").length\n\t\t\tif (numTagInputs >= 3) return\n\t\t\tconst newTagContainer = document.createElement(\"div\")\n\t\t\tnewTagContainer.className=\"tag-input-container\"\n\t\t\tconst newInput = document.createElement(\"input\")\n\t\t\tnewInput.name=`tag-${numTagInputs + 1}`\n\t\t\tnewInput.className=\"tag-input input-field\"\n\t\t\tnewInput.placeholder = \"jazz\"\n\t\t\tnewInput.type = \"text\"\n\t\t\tnewTagContainer.appendChild(newInput)\n\t\t\tconst newRemoveButton = document.createElement(\"button\")\n\t\t\tnewRemoveButton.onclick = (e) => e.target.closest(\".tag-input-container\").remove()\n\t\t\tnewRemoveButton.className=\"remove-tag-button\"\n\t\t\tconst removeIcon = document.createElement(\"div\")\n\t\t\tremoveIcon.className=\"remove-icon\"\n\t\t\tnewRemoveButton.appendChild(removeIcon)\n\t\t\tnewTagContainer.appendChild(newRemoveButton)\n\t\t\tcontainer.insertBefore(newTagContainer, elt)\n\t\t\thtmx.process(container)\n\t\t}\n\n\tasync function InitUpload(file, type) {\n\t        const url = `/init-upload/${type}`\n\t\tconst body = {\n\t\t\t'filename': file.name,\n\t\t\t'totalsize': file.size,\n\t\t\t}\n\t\n\t\tif (type == \"stems\") {\n\t\t\tbody.mainfile = htmx.find(\"input[name=audio_file_name]\").value \n\t\t}\n\t\ttry {\n\t\t      const response = await fetch(url, {\n\t\t\tmethod: 'POST',\n\t\t\tbody: JSON.stringify(body),\n\t\t      })\n\t\t      if (!response.ok) {\n\t\t\tthrow new Error(\"failed init upload\");\n\t\t      }\n\t\t      const res = await response.text()\n\t\t      message = res\n\t\t      return true\n\n\t\t    } catch (e) {\n\t\t      console.error(e)\n\t\t      return false\n\t\t    } finally {\n\t\t\tconst worker = new Worker(\"/workers/upload-worker.js\")\n\t\t\twindow.workers.push(worker)\n\t\t\tlet workerIndex = window.workers.length - 1\n\t\t\tconst fileObj = {\n\t\t\t\tfile: file,\n\t\t\t\ttype: type,\n\t\t\t} \n\t\t\twindow.workers[workerIndex].postMessage(fileObj)\n\t\t    }\n\t\t}\n\n\t\tfunction HandleArtworkSelected(elt) {\n\t\t\tconst file = elt.files[0]\n\t\t\tconst fileReader = new FileReader()\n\t\t\tfileReader.readAsDataURL(file)\n\t\t\tconst artworkInputContainer = htmx.closest(elt, \".artwork-input-container\")\n\t\t\tconst imageElt = htmx.find(artworkInputContainer, \".file-image\")\n\t\t\thtmx.on(fileReader, \"load\", () => {\n\t\t\t\timageElt.src = fileReader.result\n\t\t\t})\n\t\t\tconst selectedFile = htmx.closest(elt, \".selected-file\")\n\t\t\tconst artworkNameInput = htmx.find(selectedFile, \"input[name=artwork_file_name]\")\n\t\t\tartworkNameInput.value = file.name\n\t\t\tInitUpload(file, \"artwork\")\n\t\t}\n\n\t    \n\t\tfunction UploadAudioFile(fileObj) {\n\t\t  const worker = new Worker(\"/workers/upload-worker.js\")\n\t\t  worker.postMessage(fileObj)\n\t\t  worker.addEventListener(\"message\", (e) => {\n\t\t\tconst filename = e.data[0]\n\t\t\tconst progress = e.data[1]\n\t\t\tconst msg = e.data[2]\n\t\t\tconst fileElt = document.getElementById(`file_${e.data[0]}`)\n\t\t\tconst progressElt = htmx.find(fileElt, \".progress-marker\")\n\t\t\tprogressElt.style.width = `${e.data[1] * 100}%`\n\n\t\t\tif (msg == \"complete\") {\n\t\t\t\tconst parent = progressElt.parentElement\n\t\t\t\tparent.style.backgroundColor = 'white'\n\t\t\t\tparent.innerHTML = 'uploaded'\n\t\t\t\tparent.style.border = ''\n\t\t\t\tparent.style.fontFamily = 'Inter'\n\t\t\t\tparent.style.color = '#706f6f'\n\t\t\t\tparent.className = 'success-message'\n\t\t\t} else if (msg.includes(\"error\")) {\n\t\t\t\thtmx.find(\"#upload-error\").innerHTML = msg\t\t\n\t\t\t\tconst parent = progressElt.parentElement\n\t\t\t\tparent.style.backgroundColor = 'white'\n\t\t\t\tparent.innerHTML = 'error'\n\t\t\t\tparent.style.border = ''\n\t\t\t\tparent.style.fontFamily = 'Inter'\n\t\t\t\tparent.style.color = '#706f6f'\n\t\t\t\tparent.className = 'success-message'\n\n\t\t\t}\n\t\t  })\n\t\t}\n\t\t\n\t\twindow.addEventListener(\"upload-audio-files\", function() {\n\t\t  const audioFiles = htmx.find(\"#file-input\").files\n\n\t          for (let i = 0; i < audioFiles.length; i++) {\n\t\t\tconst filename = audioFiles[i].name\n\t\t\tconst audioFile = {file: audioFiles[i], type: \"audio\"}\n\t\t\tUploadAudioFile(audioFile)\n\t\t   }\n\t\t})\n\n\t\tfunction addStemFile(elt) {\n\t\t\tconsole.log(\"files: \", elt.files)\n\t\t\tconst files = elt.files\n\t\t\tconst stemsSection = htmx.find(\"#stems-input-section\")\n\t\t\tconst container = htmx.find(stemsSection, \".input-container\")\n\t\t\tconst addButton = htmx.find(stemsSection, \"#add-stems-button\")\n\n\t\t\tconst stemFileNamesInput = htmx.find(\"input[name=stem_file_names]\") \n\t\t\tconst stemFileNames = JSON.parse(stemFileNamesInput.value)\n\t\t\tconsole.log(stemFileNames)\n\n\t\t\tfor (let i = 0; i < files.length; i++) {\n\t\t\t\tconst fileElt = document.createElement('div');\n\t\t\t\tfileElt.innerHTML = files[i].name\n\t\t\t\tstemFileNames.push(files[i].name)\n\t\t\t\tfileElt.className = \"stem-file\"\n\t\t\t\tcontainer.insertBefore(fileElt, addButton)\n\t\t\t\tInitUpload(files[i], \"stems\")\n\t\t\t}\n\t\t\tconst newValue = JSON.stringify(stemFileNames)\n\t\t\tconsole.log(newValue)\n\t\t\tstemFileNamesInput.value = newValue\n\t\t}\n  \n\t\t</script><style id=\"upload-styles\">\n#selected-files {\n\t background-color: white;\n\t padding: 35px;\n\t border-radius: 11px;\n\t display: flex;\n}\n\n.selected-file {\n  display: flex;\n  flex-direction: row;\n  min-height: 250px;\n}\n\nlabel[for=artwork] {\n  cursor: pointer;\n  text-align: center;\n}\n\n#upload-error {\n\ttext-align: center;\n\tpadding: 10px;\n}\n\n.file-title {\n  margin-top: 10px;\n  margin-bottom: 10px;\n}\n\ninput[type=file] {\n opacity: 0%;\n width: 0;\n height: 0;\n }\n\n .file-image {\n\twidth: 100%;\n\theight: 100%;\n\tposition: absolute;\n\tcursor: pointer;\n }\n\n .input-field {\n    padding: 4px;\n    border-radius: 7px;\n    border: 1px solid #8B8E98;\n    filter: drop-shadow(0px 1px 0px #efefef) drop-shadow(0px 1px 0.5px rgba(239, 239, 239, 0.5));\n }\n\n.selected-file-left {\n  width: 200px;\n  position: relative;\n} \n\n.selected-file-right {\n  display: flex;\n  flex-direction: column;\n  justify-content: space-between; \n  padding-left: 20px; \n}\n\n\n.upload-input-form {\n  display: flex;\n  flex-direction: row;\n  gap: 15px;\n  padding-bottom: 20px;\n}\n\n.upload-input-section {\n  flex-grow: 1;\n}\n\n.input-container {\n  display: flex;\n  flex-direction: column;\n  padding-bottom: 10px;\n  position: relative;\n}\n\n.artwork-input-container {\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  width: 100%;\n  height: auto;\n  aspect-ratio: 1;\n  background-color: lightgray;\n  position: relative;\n}\n\n.artwork-input-container:hover {\n  background-color: gainsboro;\n}\n\n.progress-bar {\n  width: 100%;\n  height: 15px;\n  border: 1px solid #e5e5e5;\n  border-radius: 4px;\n}\n\n.progress-marker {\n  height: 100%;\n  background-color: var(--dark-gray);\n  width: 0;\n  border-radius: 4px;\n  }\n\n.upload-file-button {\n  background-color: var(--dark-gray);\n  text-align: center;\n  color: white;\n  border-radius: 7px;\n  cursor: pointer;\n  font-weight: 600;\n  height: 40px;\n  outline: none;\n  transition: all 0.3s cubic-bezier(0.15, 0.83, 0.66, 1);\n} \n\n\n.upload-file-button:hover {\n    background-color: white;\n    color: var(--dark-gray);\n    outline: 1px inset var(--dark-gray);\n}\n\n\n.input-form-button {\n   border:1px solid var(--dark-gray);\n   border-radius: 4px;\n   height: 20px;\n   font-size: 0.75em;\n   text-align: center;\n   margin-top: 10px;\n   padding: 4px;\n  transition: background-color 0.3s cubic-bezier(0.15, 0.83, 0.66, 1);\n}\n\n\n.input-form-button:hover { \n  color: white;\n  background-color: #2e2d2d;\n}\n\n.remove-tag-button {\n\tpadding: 4px;\n\tborder-radius: 100%;\n\tborder: 1px solid var(--dark-gray);\n\theight: 10px;\n\twidth: 10px;\n\tdisplay: flex;\n\tjustify-content: center;\n\talign-items: center;\n}\n\n\n.remove-tag-button:hover {\n\tbackground-color: #d4cfcf;\n}\n\n.tag-input-container {\n\tdisplay: flex;\n\tgap: 4px;\n\talign-items: center;\n\tmargin-bottom: 4px;\n}\n\n\n.remove-icon {\n\theight: 3px;\n\twidth: 6px;\n\tbackground-color: #8B8E98;\n\tpointer-events: auto;\n}\n\n#stems-input-section {\n\tmin-width: 150px;\n}\n\n.stem-file {\n   border:1px solid var(--dark-gray);\n   border-radius: 4px;\n   height: 20px;\n   font-size: 0.75em;\n   margin-top: 10px;\n   padding: 4px;\n   display: flex;\n   justify-content: center;\n   align-items: center;\n   background-color: #e3e3e3;\n}\n\n\t\t\t</style></head>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = uploadFormFileOnceHandle.Once().Render(templ.WithChildren(ctx, templ_7745c5c3_Var10), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

*/

var _ = templruntime.GeneratedTemplate
