"use strict";

angular.module('yobs').controller('ProfileCtrl', function($scope, facebook) {
  $scope.user = facebook.user;
});
