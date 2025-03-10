package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller"
	st "k8s.io/kubernetes/pkg/scheduler/testing"

	"github.com/muffin-rice/scheduler-plugins/pkg/apis/scheduling/v1alpha1"
	pgfake "github.com/muffin-rice/scheduler-plugins/pkg/generated/clientset/versioned/fake"
	schedinformer "github.com/muffin-rice/scheduler-plugins/pkg/generated/informers/externalversions"
	"github.com/muffin-rice/scheduler-plugins/pkg/util"
)

func Test_Run(t *testing.T) {
	ctx := context.TODO()
	createTime := metav1.Time{Time: time.Now().Add(-72 * time.Hour)}
	cases := []struct {
		name               string
		pgName             string
		minMember          uint32
		podNames           []string
		podNextPhase       v1.PodPhase
		podPhase           v1.PodPhase
		previousPhase      v1alpha1.PodGroupPhase
		desiredGroupPhase  v1alpha1.PodGroupPhase
		podGroupCreateTime *metav1.Time
	}{
		{
			name:              "Group running",
			pgName:            "pg1",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodRunning,
			previousPhase:     v1alpha1.PodGroupScheduled,
			desiredGroupPhase: v1alpha1.PodGroupRunning,
		},
		{
			name:              "Group running, more than min member",
			pgName:            "pg11",
			minMember:         2,
			podNames:          []string{"pod11", "pod21"},
			podPhase:          v1.PodRunning,
			previousPhase:     v1alpha1.PodGroupScheduled,
			desiredGroupPhase: v1alpha1.PodGroupRunning,
		},
		{
			name:              "Group failed",
			pgName:            "pg2",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodFailed,
			previousPhase:     v1alpha1.PodGroupScheduled,
			desiredGroupPhase: v1alpha1.PodGroupFailed,
		},
		{
			name:              "Group finished",
			pgName:            "pg3",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodSucceeded,
			previousPhase:     v1alpha1.PodGroupScheduled,
			desiredGroupPhase: v1alpha1.PodGroupFinished,
		},
		{
			name:              "Group status convert from scheduling to scheduled",
			pgName:            "pg4",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodPending,
			previousPhase:     v1alpha1.PodGroupScheduling,
			desiredGroupPhase: v1alpha1.PodGroupScheduled,
		},
		{
			name:              "Group status convert from scheduling to succeed",
			pgName:            "pg5",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodPending,
			previousPhase:     v1alpha1.PodGroupScheduling,
			desiredGroupPhase: v1alpha1.PodGroupFinished,
			podNextPhase:      v1.PodSucceeded,
		},
		{
			name:              "Group status convert from scheduling to succeed",
			pgName:            "pg6",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodPending,
			previousPhase:     v1alpha1.PodGroupScheduling,
			desiredGroupPhase: v1alpha1.PodGroupFinished,
			podNextPhase:      v1.PodSucceeded,
		},
		{
			name:              "Group status convert from pending to prescheduling",
			pgName:            "pg7",
			minMember:         2,
			podNames:          []string{"pod1", "pod2"},
			podPhase:          v1.PodPending,
			previousPhase:     v1alpha1.PodGroupPending,
			desiredGroupPhase: v1alpha1.PodGroupFinished,
			podNextPhase:      v1.PodSucceeded,
		},
		{
			name:               "Group group should not enqueue, created too long",
			pgName:             "pg8",
			minMember:          2,
			podNames:           []string{"pod1", "pod2"},
			podPhase:           v1.PodRunning,
			previousPhase:      v1alpha1.PodGroupPending,
			desiredGroupPhase:  v1alpha1.PodGroupPending,
			podGroupCreateTime: &createTime,
		},
		{
			name:              "Group group min member more than Pod number",
			pgName:            "pg9",
			minMember:         3,
			podNames:          []string{"pod91", "pod92"},
			podPhase:          v1.PodPending,
			previousPhase:     v1alpha1.PodGroupPending,
			desiredGroupPhase: v1alpha1.PodGroupPreScheduling,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ps := makePods(c.podNames, c.pgName, c.podPhase)
			kubeClient := fake.NewSimpleClientset(ps[0], ps[1])
			pg := makePG(c.pgName, 2, c.previousPhase, c.podGroupCreateTime)
			pgClient := pgfake.NewSimpleClientset(pg)

			informerFactory := informers.NewSharedInformerFactory(kubeClient, controller.NoResyncPeriodFunc())
			pgInformerFactory := schedinformer.NewSharedInformerFactory(pgClient, controller.NoResyncPeriodFunc())
			podInformer := informerFactory.Core().V1().Pods()
			pgInformer := pgInformerFactory.Scheduling().V1alpha1().PodGroups()
			ctrl := NewPodGroupController(kubeClient, pgInformer, podInformer, pgClient)

			pgInformerFactory.Start(ctx.Done())
			informerFactory.Start(ctx.Done())
			// 0 means not set
			if len(c.podNextPhase) != 0 {
				ps := makePods(c.podNames, c.pgName, c.podNextPhase)
				for _, p := range ps {
					kubeClient.CoreV1().Pods(p.Namespace).UpdateStatus(ctx, p, metav1.UpdateOptions{})
				}
			}
			go ctrl.Run(1, ctx.Done())
			err := wait.Poll(200*time.Millisecond, 1*time.Second, func() (done bool, err error) {
				pg, err := pgClient.SchedulingV1alpha1().PodGroups("default").Get(ctx, c.pgName, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				if pg.Status.Phase != c.desiredGroupPhase {
					return false, fmt.Errorf("want %v, got %v", c.desiredGroupPhase, pg.Status.Phase)
				}
				return true, nil
			})
			if err != nil {
				t.Fatal("Unexpected error", err)
			}
		})
	}

}

func makePods(podNames []string, pgName string, phase v1.PodPhase) []*v1.Pod {
	pds := make([]*v1.Pod, 0)
	for _, name := range podNames {
		pod := st.MakePod().Namespace("default").Name(name).Obj()
		pod.Labels = map[string]string{util.PodGroupLabel: pgName}
		pod.Status.Phase = phase
		pds = append(pds, pod)
	}
	return pds
}

func makePG(pgName string, minMember int32, previousPhase v1alpha1.PodGroupPhase, createTime *metav1.Time) *v1alpha1.PodGroup {
	pg := &v1alpha1.PodGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:              pgName,
			Namespace:         "default",
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Spec: v1alpha1.PodGroupSpec{
			MinMember: minMember,
		},
		Status: v1alpha1.PodGroupStatus{
			OccupiedBy:        "test",
			Scheduled:         minMember,
			ScheduleStartTime: metav1.Time{Time: time.Now()},
			Phase:             previousPhase,
		},
	}
	if createTime != nil {
		pg.CreationTimestamp = *createTime
	}
	return pg
}
