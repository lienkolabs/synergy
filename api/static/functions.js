// console.log("estou aqui")

window.onload = function () {

  // general

  let reasonsf = document.getElementById("reasonsfield");
  if (reasonsf) {    
    reasonsf.addEventListener("focusin",displayinfo("reasonsfieldinfo"));
    reasonsf.addEventListener("focusout",hideinfo("reasonsfieldinfo"));
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

  // update collective

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

  // update board

  let newdescrboard = document.getElementById("newdescriptionboard");
  if (newdescrboard) {    
    newdescrboard.addEventListener("focusin",displayinfo("newdescriptionboardinfo"));
    newdescrboard.addEventListener("focusout",hideinfo("newdescriptionboardinfo"));
  }

  let newkeywsboard = document.getElementById("newkeywordsboard");
  if (newkeywsboard) {    
    newkeywsboard.addEventListener("focusin",displayinfo("newkeywordsboardinfo"));
    newkeywsboard.addEventListener("focusout",hideinfo("newkeywordsboardinfo"));
  }

  let newpmajboard = document.getElementById("newpinmajboard");
  if (newpmajboard) {    
    newpmajboard.addEventListener("focusin",displayinfo("newpinmajboardinfo"));
    newpmajboard.addEventListener("focusout",hideinfo("newpinmajboardinfo"));
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

  // update event 

  let newdescrevent = document.getElementById("newdescriptionevent");
  if (newdescrevent) {    
    newdescrevent.addEventListener("focusin",displayinfo("descriptioneventinfo"));
    newdescrevent.addEventListener("focusout",hideinfo("descriptioneventinfo"));
  }

  let newvenueevent = document.getElementById("newvenueevent");
  if (newvenueevent) {    
    newvenueevent.addEventListener("focusin",displayinfo("newvenueeventinfo"));
    newvenueevent.addEventListener("focusout",hideinfo("newvenueeventinfo"));
  }

  let newopen = document.getElementById("newopen");
  if (newopen) {    
    newopen.addEventListener("focusin",displayinfo("newopeneventinfo"));
    newopen.addEventListener("focusout",hideinfo("newopeneventinfo"));
  }

  let newpublic = document.getElementById("newpublic");
  if (newpublic) {    
    newpublic.addEventListener("focusin",displayinfo("newpubliceventinfo"));
    newpublic.addEventListener("focusout",hideinfo("newpubliceventinfo"));
  }

  let newpevent = document.getElementById("newpmanagementevent");
  if (newpevent) {    
    newpevent.addEventListener("focusin",displayinfo("newpmanagementeventinfo"));
    newpevent.addEventListener("focusout",hideinfo("newpmanagementeventinfo"));
  }

  // new draft

  let singledraft = document.getElementById("single");
  if (singledraft) {    
    singledraft.addEventListener("focusin",displayinfo("singledraftinfo"));
    singledraft.addEventListener("focusout",hideinfo("singledraftinfo"));
  }

  let coautdraft = document.getElementById("coauthored");
  if (coautdraft) {    
    coautdraft.addEventListener("focusin",displayinfo("coauthoreddraftinfo"));
    coautdraft.addEventListener("focusout",hideinfo("coauthoreddraftinfo"));
  }

  let colobodraft = document.getElementById("collectively");
  if (colobodraft) {    
    colobodraft.addEventListener("focusin",displayinfo("collectivelydraftinfo"));
    colobodraft.addEventListener("focusout",hideinfo("collectivelydraftinfo"));
  }

  let coauthorsdraft = document.getElementById("coauthorsdraft");
  if (coauthorsdraft) {    
    coauthorsdraft.addEventListener("focusin",displayinfo("coauthorsdraftinfo"));
    coauthorsdraft.addEventListener("focusout",hideinfo("coauthorsdraftinfo"));
  }

  let policydraft = document.getElementById("policydraft");
  if (policydraft) {    
    policydraft.addEventListener("focusin",displayinfo("policydraftinfo"));
    policydraft.addEventListener("focusout",hideinfo("policydraftinfo"));
  }

  let collectivedraft = document.getElementById("collectivedraft");
  if (collectivedraft) {    
    collectivedraft.addEventListener("focusin",displayinfo("collectivedraftinfo"));
    collectivedraft.addEventListener("focusout",hideinfo("collectivedraftinfo"));
  }

  let titledraft = document.getElementById("titledraft");
  if (titledraft) {    
    titledraft.addEventListener("focusin",displayinfo("titledraftinfo"));
    titledraft.addEventListener("focusout",hideinfo("titledraftinfo"));
  }

  let keywdraft = document.getElementById("keywordsdraft");
  if (keywdraft) {    
    keywdraft.addEventListener("focusin",displayinfo("keywordsdraftinfo"));
    keywdraft.addEventListener("focusout",hideinfo("keywordsdraftinfo"));
  }

  let descrdraft = document.getElementById("descriptiondraft");
  if (descrdraft) {    
    descrdraft.addEventListener("focusin",displayinfo("descriptiondraftinfo"));
    descrdraft.addEventListener("focusout",hideinfo("descriptiondraftinfo"));
  }

  let fileudraft = document.getElementById("fileudraft");
  if (fileudraft) {    
    fileudraft.addEventListener("focusin",displayinfo("fileudraftinfo"));
    fileudraft.addEventListener("focusout",hideinfo("fileudraftinfo"));
  }

  let previousvdraft = document.getElementById("previousvdraft");
  if (previousvdraft) {    
    previousvdraft.addEventListener("focusin",displayinfo("previousvdraftinfo"));
    previousvdraft.addEventListener("focusout",hideinfo("previousvdraftinfo"));
  }

  let referencesdraft = document.getElementById("referencesdraft");
  if (referencesdraft) {    
    referencesdraft.addEventListener("focusin",displayinfo("referencesdraftinfo"));
    referencesdraft.addEventListener("focusout",hideinfo("referencesdraftinfo"));
  }

  // modal 
  let updatebutton = document.getElementById("instructionmodal");
  let cancelbutton = document.getElementById("cancel");
  let dialog = document.getElementById("dialogreasons");
  console.log(dialog)
  
  if (dialog) {
    dialogteste(dialog, updatebutton, cancelbutton);
  }
}

// form info functions

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

// modal functions

function dialogteste(dialog, updatebutton, cancelbutton) {
  dialog.returnValue = "teste";

    // Update button opens a modal dialog
    updatebutton.addEventListener("click", () => {
      dialog.showModal();
      openCheck(dialog);
    });
  
    // Form cancel button closes the dialog box
    cancelbutton.addEventListener("click", () => {
      dialog.close("animalNotChosen");
      openCheck(dialog);
    });
}

function openCheck(dialog) {
  if (dialog.open) {
    console.log("Dialog open");
  } else {
    console.log("Dialog closed");
  }
}