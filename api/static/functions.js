console.log("estou aqui")

window.onload = function () {

  // general

  let reasonsf = document.getElementById("reasonsfield");
  if (reasonsf) {    
    reasonsf.addEventListener("focusin",displayinfo("reasonsfieldinfo"));
    reasonsf.addEventListener("focusout",hideinfo("reasonsfieldinfo"));
  }

  let publicevent = document.getElementById("public");
  if (publicevent) {    
    publicevent.addEventListener("focusin",displayinfo("publiceventinfo"));
    publicevent.addEventListener("focusout",hideinfo("publiceventinfo"));
  }

  let openevent = document.getElementById("open");
  if (openevent) {    
    openevent.addEventListener("focusin",displayinfo("openeventinfo"));
    openevent.addEventListener("focusout",hideinfo("openeventinfo"));
  }

  // create collective
  
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

// create board

  let nameboard = document.getElementById("nameboard");
  if (nameboard) {    
    nameboard.addEventListener("focusin",displayinfo("nameboardinfo"));
    nameboard.addEventListener("focusout",hideinfo("nameboardinfo"));
  }

  let descrboard = document.getElementById("descriptionboard");
  if (descrboard) {    
    descrboard.addEventListener("focusin",displayinfo("descriptionboardinfo"));
    descrboard.addEventListener("focusout",hideinfo("descriptionboardinfo"));
  }

  let keywboard = document.getElementById("keywordsboard");
  if (keywboard) {    
    keywboard.addEventListener("focusin",displayinfo("keywordsboardinfo"));
    keywboard.addEventListener("focusout",hideinfo("keywordsboardinfo"));
  }

  let pinpolboard = document.getElementById("pinpolicyboard");
  if (pinpolboard) {    
    pinpolboard.addEventListener("focusin",displayinfo("pinpolicyboardinfo"));
    pinpolboard.addEventListener("focusout",hideinfo("pinpolicyboardinfo"));
  }

  // create event 

  let descrevent = document.getElementById("descriptionevent");
  if (descrevent) {    
    descrevent.addEventListener("focusin",displayinfo("descriptioneventinfo"));
    descrevent.addEventListener("focusout",hideinfo("descriptioneventinfo"));
  }

  let startevent = document.getElementById("startatevent");
  if (startevent) {    
    startevent.addEventListener("focusin",displayinfo("startateventinfo"));
    startevent.addEventListener("focusout",hideinfo("startateventinfo"));
  }

  let endevent = document.getElementById("estimatedendevent");
  if (endevent) {    
    endevent.addEventListener("focusin",displayinfo("estimatedeventinfo"));
    endevent.addEventListener("focusout",hideinfo("estimatedeventinfo"));
  }

  let venevent = document.getElementById("venueevent");
  if (venevent) {    
    venevent.addEventListener("focusin",displayinfo("venueeventinfo"));
    venevent.addEventListener("focusout",hideinfo("venueeventinfo"));
  }

  let polmanevent = document.getElementById("policymanagementevent");
  if (polmanevent) {    
    polmanevent.addEventListener("focusin",displayinfo("policymanagementeventinfo"));
    polmanevent.addEventListener("focusout",hideinfo("policymanagementeventinfo"));
  }

  let mangevent = document.getElementById("managersevent");
  if (mangevent) {    
    mangevent.addEventListener("focusin",displayinfo("managerseventinfo"));
    mangevent.addEventListener("focusout",hideinfo("managerseventinfo"));
  }

  // update event

  let newdescrcol = document.getElementById("newdescriptioncollective");
  if (newdescrcol) {    
    newdescrcol.addEventListener("focusin",displayinfo("newdescriptioncollectiveinfo"));
    newdescrcol.addEventListener("focusout",hideinfo("newdescriptioncollectiveinfo"));
  }

  let newpolmcol = document.getElementById("newpolicymajcollective");
  if (newpolmcol) {    
    newpolmcol.addEventListener("focusin",displayinfo("newpolicymajcollectiveinfo"));
    newpolmcol.addEventListener("focusout",hideinfo("newpolicymajcollectiveinfo"));
  }

  let newsuppolmcol = document.getElementById("newpolicysupermajcollective");
  if (newsuppolmcol) {    
    newsuppolmcol.addEventListener("focusin",displayinfo("newpolicysupermajcollectiveinfo"));
    newsuppolmcol.addEventListener("focusout",hideinfo("newpolicysupermajcollectiveinfo"));
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

