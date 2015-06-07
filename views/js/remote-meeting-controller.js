/**
 * Created by rang on 2015-06-07.
 */
angular.module('remotemeeting.controllers', [])
    .controller('RemoteMeetingController', function($scope, $http) {
        var localStream;
        var localPeerConnection;
        var remotePeerConnection;

        var localVideo = document.getElementById('localVideo');
        var remoteVideo = document.getElementById('remoteVideo');

        localVideo.addEventListener('loadedmetadata', function() {
            trace('Local video currentSrc: ' + this.currentSrc +
                ', videoWidth: ' + this.videoWidth +
                'px,  videoHeight: ' + this.videoHeight + 'px');
        });

        remoteVideo.addEventListener('loadedmetadata', function() {
            trace('Remote video currentSrc: ' + this.currentSrc +
                ', videoWidth: ' + this.videoWidth +
                'px,  videoHeight: ' + this.videoHeight + 'px');
        });

        $scope.startButtonDisabled = false;
        $scope.callButtonDisabled = true;
        $scope.hangupButtonDisabled = true;

        var total = '';

        function trace(text) {
            total += text;
            console.log((window.performance.now() / 1000).toFixed(3) + ': ' + text);
        }

        $scope.start = function() {
            $scope.startButtonDisabled = true;
            navigator.getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia;
            navigator.getUserMedia(
                {
                    video: true
                },
                function (stream) {
                    localVideo.src = URL.createObjectURL(stream);
                    localStream = stream;
                    $scope.callButtonDisabled = false;
                },
                function(error) {
                    trace('navigator.getUserMedia error: ', error);
                }
            );
        };

        $scope.call = function() {
            $scope.callButtonDisabled = true;
            $scope.hangupButtonDisabled = false;

            if (localStream.getVideoTracks().length > 0) {
                trace('Using video device: ' + localStream.getVideoTracks()[0].label);
            }
            if (localStream.getAudioTracks().length > 0) {
                trace('Using audio device: ' + localStream.getAudioTracks()[0].label);
            }
            var servers = null;
            localPeerConnection = new webkitRTCPeerConnection(servers);
            localPeerConnection.onicecandidate = function (event) {
                if (event.candidate) {
                    remotePeerConnection.addIceCandidate(new RTCIceCandidate(event.candidate));
                    trace('Local ICE candidate: \n' + event.candidate.candidate);
                }
            };

            remotePeerConnection = new webkitRTCPeerConnection(servers);
            remotePeerConnection.onicecandidate = function (event) {
                if (event.candidate) {
                    localPeerConnection.addIceCandidate(new RTCIceCandidate(event.candidate));
                    trace('Remote ICE candidate: \n ' + event.candidate.candidate);
                }
            };
            remotePeerConnection.onaddstream = function (event) {
                remoteVideo.src = URL.createObjectURL(event.stream);
                trace('Received remote stream');
            };

            localPeerConnection.addStream(localStream);
            localPeerConnection.createOffer(function (description) {
                localPeerConnection.setLocalDescription(description);
                trace('Offer from localPeerConnection: \n' + description.sdp);
                remotePeerConnection.setRemoteDescription(description);
                remotePeerConnection.createAnswer(function (description) {
                    remotePeerConnection.setLocalDescription(description);
                    trace('Answer from remotePeerConnection: \n' + description.sdp);
                    localPeerConnection.setRemoteDescription(description);
                });
            });
        };

        $scope.hangup = function() {
            trace('Ending call');
            localPeerConnection.close();
            remotePeerConnection.close();
            localPeerConnection = null;
            remotePeerConnection = null;
            $scope.hangupButtonDisabled = true;
            $scope.callButtonDisabled = false;
        };

    });
