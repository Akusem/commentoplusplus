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
    $("#add-new-label").hide()
    $("#label-creator").show()
  }

  global.closeLabelCreator = function() {
    $("#label-creator").hide()
    $("#add-new-label").show()
  }

  global.setAllLabels = function() {
    // Get the list of labels already created if labels are allowed
    var data = global.dashboard.$data;

    if (data.domains[data.cd].allowLabels) {
      var json = {
        "ownerToken": global.cookieGet("commentoOwnerToken"),
        "domain": data.domains[data.cd].domain
      };
      global.post(global.origin + "/api/label/owner/listAll", json, function(resp) {
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

} (window.commento, document));
