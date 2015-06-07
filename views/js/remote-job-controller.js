/**
 * Created by rang on 2015-06-07.
 */
angular.module('remotejob.controllers', [])
    .controller('RemoteJobController', function($scope, $http) {
        $scope.getList = function() {
            $http.get('/api/getRemoteJobInfo').
                success(function(data, status, headers, config) {
                    $scope.list = data.ResultContent
                    console.log(data.ResultContent);
                }).
                error(function(data, status, headers, config) {
                });
        };

        $scope.startBatch = function(){
            $http.get('/api/batchNow').
                success(function(data, status, headers, config) {
                    if(data.ResultCode == 200) {
                        alert("bach complete");
                        $scope.getList();
                    }
                }).
                error(function(data, status, headers, config) {
                });
        };
        $scope.getList();
    });