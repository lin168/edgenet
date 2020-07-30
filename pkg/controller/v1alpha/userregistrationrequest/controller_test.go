package userregistrationrequest

import (
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStartController(t *testing.T) {
	g := URRTestGroup{}
	g.Init()
	// Run the controller in a goroutine
	go Start(g.client, g.edgenetclient)
	// Create a user registration object
	g.edgenetclient.AppsV1alpha().UserRegistrationRequests(fmt.Sprintf("authority-%s", g.authorityObj.GetName())).Create(g.userRegistrationObj.DeepCopy())
	// Wait for the status update of created object
	time.Sleep(time.Millisecond * 500)
	// Get the object and check the status
	AR, _ := g.edgenetclient.AppsV1alpha().UserRegistrationRequests(fmt.Sprintf("authority-%s", g.authorityObj.GetName())).Get(g.userRegistrationObj.GetName(), metav1.GetOptions{})
	if AR.Status.Expires == nil || AR.Status.Message == nil {
		t.Error(ErrorDict["add-func"])
	}
	// Update a Authority request
	// Update contact email
	g.userRegistrationObj.Spec.Email = "URR@edge-net.org"
	g.userRegistrationObj.Spec.Approved = true
	g.edgenetclient.AppsV1alpha().UserRegistrationRequests(fmt.Sprintf("authority-%s", g.authorityObj.GetName())).Update(g.userRegistrationObj.DeepCopy())
	time.Sleep(time.Millisecond * 500)
	// Checking if user registration transitioned to user after update
	user, _ := g.edgenetclient.AppsV1alpha().Users(fmt.Sprintf("authority-%s", g.authorityObj.GetName())).Get(g.userRegistrationObj.GetName(), metav1.GetOptions{})
	if user == nil {
		t.Error(ErrorDict["upd-func"])
	}
}
