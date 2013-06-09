"use strict";

angular.module('yobs').factory('facebook', ['$window', '$rootScope', function($window, $rootScope) {
    var user = {};
    $window.fbAsyncInit = function() {
      FB.init({
        appId      : '470078866407798', // App ID
        channelUrl : '//' + this.location.host + '/fb_channel.html', // Channel File
        status     : true, // check login status
        cookie     : true, // enable cookies to allow the server to access the session
        xfbml      : true  // parse XFBML
      });

      // Here we subscribe to the auth.authResponseChange JavaScript event.
      // This event is fired for any authentication related change, such as login,
      // logout or session refresh. This means that whenever someone who was
      // previously logged out tries to log in again, the correct case below will
      // be handled. 
      FB.Event.subscribe('auth.authResponseChange', function(response) {
        // Here we specify what we do with the response anytime this event occurs. 
        if (response.status === 'connected') {
          // The response object is returned with a status field that lets the app know the current
          // login status of the person. In this case, we're handling the situation where they 
          // have logged in to the app.
          FB.api('/me', function(response) {
            $rootScope.$apply(function() {
              user = response;
            });
          });
        } else if (response.status === 'not_authorized') {
          // In this case, the person is logged into Facebook, but not into the app,
        } else {
          // In this case, the person is not logged into Facebook.
        }
      });
    };

    // Load the Facebook SDK asynchronously
    (function(d){
       var js, id = 'facebook-jssdk', ref = d.getElementsByTagName('script')[0];
       if (d.getElementById(id)) {return;}
       js = d.createElement('script'); js.id = id; js.async = true;
       js.src = "//connect.facebook.net/en_US/all.js";
       ref.parentNode.insertBefore(js, ref);
     }($window.document));

    return {
      user: function() {
        return user;
      }
    };
}]);

//angular.injector([fbServiceModule]).get('facebook');
