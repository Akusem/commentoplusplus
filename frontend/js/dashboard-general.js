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
  }

  global.closeLabelCreator = function() {
    $("#label-creator").hide()
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

} (window.commento, document));
