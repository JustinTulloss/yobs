"use strict";

angular.module('yobs').controller('LoginCtrl', function($scope, $timeout, facebook) {
  facebook.scope.$watch('status', function(status) {
    $scope.fbStatus = status;
    $timeout(function() {
      facebook.parseXFBML();
    }, 1);
  });
  facebook.scope.$watch('user', function(user) {
    $scope.name = user ? user.name : "";
  });
});
