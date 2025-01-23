
  const RemoveAccountDropdown = (e) => {
    if (e.relatedTarget == null) {
      document.getElementById("account-dropdown").innerHTML = ""
      document.getElementById("account-dropdown").style.visibility = "hidden";
    } else if (e.relatedTarget.id != "profile-photo-button") {
      document.getElementById("account-dropdown").innerHTML = ""
      document.getElementById("account-dropdown").style.visibility = "hidden";
    }
    return
  }

  const FocusAccountDropdown = () => {
    console.log("focus account dropdown")
    let dropdown = document.getElementById("account-dropdown")
    dropdown.focus()
  }
