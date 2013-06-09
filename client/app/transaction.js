"use strict";

angular.module('yobs').controller('NewTxnCtrl', function($scope, facebook)  {
  $scope.friends = [];
  facebook.scope.$watch('status', function(status) {
    $scope.fbStatus = status;
    if (status === 'connected') {
      var def = facebook.friends(['name', 'picture']);
      def.then(function(friends) {
        $scope.friends = friends;
      }, function(error) {
        console.error(error.message, error);
      });
    }
  });
});
