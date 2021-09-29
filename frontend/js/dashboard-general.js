(function (global, document) {
  "use strict";

  (document);

  // Opens the general settings window.
  global.generalOpen = function() {
    $(".view").hide();
    $("#general-view").show();

    global.setAllLabels();
  };

  global.generalSaveHandler = function() {
    var data = global.dashboard.$data;

    global.buttonDisable("#save-general-button");
    global.domainUpdate(data.domains[data.cd], function() {
      global.globalOKShow("Settings saved!");
      global.buttonEnable("#save-general-button");
    });
  };

  global.ssoProviderChangeHandler = function() {
    var data = global.dashboard.$data;

    if (data.domains[data.cd].ssoSecret === "") {
      var json = {
        "ownerToken": global.cookieGet("commentoOwnerToken"),
        "domain": data.domains[data.cd].domain,
      };

      global.post(global.origin + "/api/domain/sso/new", json, function(resp) {
        if (!resp.success) {
          global.globalErrorShow(resp.message);
          return;
        }

        data.domains[data.cd].ssoSecret = resp.ssoSecret;
        $("#sso-secret").val(data.domains[data.cd].ssoSecret);
      });
    } else {
      $("#sso-secret").val(data.domains[data.cd].ssoSecret);
    }
  };

  global.openLabelCreator = function() {
    $("#label-creator").show()
    // Avoid having label Creator & Editor open at the same time
    global.closeLabelEditor()
  }

  global.closeLabelCreator = function() {
    $("#label-creator").hide()
  }

  global.openLabelEditor = function(label) {
    var data = global.dashboard.$data
    // Set label Editor to label of interest value
    data.domains[data.cd].editLabelHex = label.labelHex;
    data.domains[data.cd].editLabelName = label.name;
    data.domains[data.cd].editLabelColor = label.color;
    $("#label-editor").show()
    // Avoid having label Creator & Editor open at the same time
    global.closeLabelCreator()
  }

  global.closeLabelEditor = function() {
    $("#label-editor").hide()
  }

  global.setAllLabels = function() {
    // Get the list of labels already created if labels are allowed
    var data = global.dashboard.$data;

    if (data.domains[data.cd].allowLabels) {
      var json = {
        "domain": data.domains[data.cd].domain
      };
      global.post(global.origin + "/api/label/listAll", json, function(resp) {
        if (!resp.success) {
          global.globalErrorShow(resp.message);
          return
        }

        Vue.set(data.domains[data.cd], "labelsAll", resp.labels)
      });
    }
  }

  global.createLabel = function() {
    var data = global.dashboard.$data;
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "name": data.domains[data.cd].newLabelName,
      "color": data.domains[data.cd].newLabelColor
    };

    if (json.name === "" || json.color === "") {
      global.globalErrorShow("Label Name missing");
      return
    }

    global.post(global.origin + "/api/label/new", json, function(resp) {
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return
      }

      global.closeLabelCreator();
      global.setAllLabels();
    });
  }

  global.editLabel = function() {
    var data = global.dashboard.$data
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "name": data.domains[data.cd].editLabelName,
      "color": data.domains[data.cd].editLabelColor,
      "labelHex": data.domains[data.cd].editLabelHex,
    };

    global.post(global.origin + "/api/label/edit", json, function(resp) {
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return
      }
      // Reload all labels, new one included
      global.closeLabelEditor();
      global.setAllLabels();
    });
  }

  global.deleteLabel = function() {
    var data = global.dashboard.$data
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "labelHex": data.domains[data.cd].labelHexToDelete,
    };

    global.post(global.origin + "/api/label/delete", json, function(resp) {
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return
      }
      // Reset label creator to default value
      data.domains[data.cd].newLabelName = "";
      data.domains[data.cd].newLabelColor = "#44ad8e";
      // Reload all labels, new one included
      global.setAllLabels();
      // Close delete modal
      document.location.hash = "#modal-close"
    });
  }

  global.openDeleteLabelModal = function(labelHex) {
    var data = global.dashboard.$data
    data.domains[data.cd].labelHexToDelete = labelHex
    document.location.hash = "#delete-label-modal"
  }

} (window.commento, document));
