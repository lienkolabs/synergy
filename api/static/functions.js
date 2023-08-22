// console.log("estou aqui")

window.onload = function () {

  // input fields detailing
  let elementsdetailed = document.getElementsByClassName("detailed");
  if (elementsdetailed) {
    for (let i=0; i<elementsdetailed.length; i++) {
      let el = elementsdetailed[i];
      let id = el.getAttribute("id");
      el.addEventListener("focusin", displayinfo(id+"info"));
      el.addEventListener("focusout", hideinfo(id+"info"));
    }
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

function closedialog(id) {
  let el = document.getElementById(id);
  el.close();
}

function dialogreasonsreact() {
  let el = document.getElementById("dialogreasons");
  el.showModal();
  let reaction = document.getElementById("reaction");
  let reactionmodal = document.getElementById("reactionmodal");
  reactionmodal.checked = reaction.checked;
}

function dialogreasonsreact() {
  // shows dialog element
  let el = document.getElementById("dialogreasons");
  el.showModal();

  // gets reaction from main page into dialog
  let reaction = document.getElementById("reaction");
  let reactionmodal = document.getElementById("reactionmodal");
  reactionmodal.checked = reaction.checked;
  
  // gets outline paragraph to be shown in modal
  let reactpar = document.getElementById("reactionoutline");
  let pagename = document.getElementById("modaloutlinename").innerHTML;

  if (reactionmodal.checked) {
    reactpar.innerHTML = "like "+ pagename;
  } else {
    reactpar.innerHTML = "dislike " + pagename;
  }
}
