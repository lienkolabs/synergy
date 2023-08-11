console.log("estou aqui")

window.onload = function () {
  let namecol = document.getElementById("namecollective");
  if (namecol) {    
    namecol.addEventListener("focusin",displayinfo("namecollectiveinfo"));
    namecol.addEventListener("focusout",hideinfo("namecollectiveinfo"));
  }

  let descrcol = document.getElementById("descriptioncollective");
  if (descrcol) {    
    descrcol.addEventListener("focusin",displayinfo("descriptioncollectiveinfo"));
    descrcol.addEventListener("focusout",hideinfo("descriptioncollectiveinfo"));
  }

  let polcol = document.getElementById("policycollective");
  if (polcol) {    
    polcol.addEventListener("focusin",displayinfo("policycollectiveinfo"));
    polcol.addEventListener("focusout",hideinfo("policycollectiveinfo"));
  }

  let spolcol = document.getElementById("superpolicycollective");
  if (spolcol) {    
    spolcol.addEventListener("focusin",displayinfo("superpolicycollectiveinfo"));
    spolcol.addEventListener("focusout",hideinfo("superpolicycollectiveinfo"));
  }

  let reasonsf = document.getElementById("reasonsfield");
  if (reasonsf) {    
    reasonsf.addEventListener("focusin",displayinfo("reasonsfieldinfo"));
    reasonsf.addEventListener("focusout",hideinfo("reasonsfieldinfo"));
  }
}

function displayinfo(id) {
  return () => {
  let el = document.getElementById(id);
  el.classList.remove("fieldinfohide");
  el.classList.add("fieldinfoshow");
  }
}

function hideinfo(id) {
	return () => {
    let el = document.getElementById(id)
	  el.classList.remove("fieldinfoshow");
  	el.classList.add("fieldinfohide");
   }
}

