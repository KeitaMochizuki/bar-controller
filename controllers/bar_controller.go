/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"           // apimachineryでruntime-objectを使う
	ctrl "sigs.k8s.io/controller-runtime"       // controller-runtimeを使う
	"sigs.k8s.io/controller-runtime/pkg/client" // controller-runtimeのclientを使う
	"sigs.k8s.io/controller-runtime/pkg/log"

	barcontrollerv1 "my.domain/bar-controller/api/v1" // Barリソースを扱う
)

// BarReconciler reconciles a Bar object
type BarReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=barcontroller.my.domain,resources=bars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=barcontroller.my.domain,resources=bars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=barcontroller.my.domain,resources=bars/finalizers,verbs=update

func (r *BarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	Log := log.FromContext(ctx)

	// [Debug]Reconcileが呼ばれたタイミングでログ出力
	Log.Info("★Reconcile Function Called★")
	Log.Info("[Debug]req(NamespacedName):" + req.Namespace + "/" + req.Name)

	// ① Barオブジェクトを取得する(In-Memory-Cacheから)
	var bar barcontrollerv1.Bar
	Log.Info("[Step1]Fetching Bar Object")
	// reconcile対象のBarオブジェクト（req.NamespacedNameとして引数で渡ってくる）を取得
	if err := r.Get(ctx, req.NamespacedName, &bar); err != nil {
		Log.Error(err, "Unable to get Bar from in-memory-cache")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	Log.Info("[Debug] Nginx: " + bar.Namespace + "/" + bar.Name + "from in-memory-cache")

	// ② bar.spac.messageとbar.status.messageを比較して差分があれば更新

	if bar.Spec.Message != bar.Status.Message {
		bar.Status.Message = "Hello " + bar.Spec.Message
		if err := r.Status().Update(ctx, &bar); err != nil {
			Log.Error(err, "Unable to update Bar status bar.status.message")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&barcontrollerv1.Bar{}).
		Complete(r)
}
