package logic

import (
	"context"
	"fmt"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func SyncDeployments(sourceCli, targetCli *kubernetes.Clientset, namespace string) error {
	sourceList, err := sourceCli.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list source deployments: %v", err)
	}

	for _, source := range sourceList.Items {
		fmt.Println(source.Name)

		target, err := targetCli.AppsV1().Deployments(namespace).Get(context.TODO(), source.Name, v1.GetOptions{})

		if err != nil {
			// 如果目标集群中不存在该 Deployment，则创建一个新的 Deployment
			if kerrors.IsNotFound(err) {
				// 清理数据
				source.ResourceVersion = ""
				_, err = targetCli.AppsV1().Deployments(namespace).Create(context.TODO(), &source, v1.CreateOptions{})
				if err != nil {
					log.Printf("创建目标 deployment 失败 %s: %v", source.Name, err)
					continue
				}
				log.Printf("在目标集群中成功创建了 deployment %s", source.Name)
			} else {
				log.Printf("获取目标 deployment 失败 %s: %v", source.Name, err)
			}
			continue
		}
		// 如果存在，则更新目标集群中的 Deployment
		//// 更新目标集群相关配置
		target.Spec.Replicas = source.Spec.Replicas
		target.Spec.Template.Spec.Containers[0].Image = source.Spec.Template.Spec.Containers[0].Image

		targetRes, err := targetCli.AppsV1().Deployments(namespace).Update(context.TODO(), target, v1.UpdateOptions{})
		if err != nil {
			log.Printf("更新目标 deployment 失败 %s: %v", targetRes.Name, err)
		} else {
			log.Printf("成功同步 deployment %s 副本数 %v", targetRes.Name, *targetRes.Spec.Replicas)
		}

	}

	return nil
}
