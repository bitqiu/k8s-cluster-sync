package logic

import (
	"context"
	"fmt"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func SyncService(sourceCli, targetCli *kubernetes.Clientset, namespace string) error {
	sourceList, err := sourceCli.CoreV1().Services(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list source services: %v", err)
	}

	for _, source := range sourceList.Items {
		targetService, err := targetCli.CoreV1().Services(namespace).Get(context.TODO(), source.Name, v1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				// 如果目标集群中不存在该 Service，则创建一个新的 Service
				source.Spec.ClusterIP = ""
				source.ResourceVersion = ""
				_, err = targetCli.CoreV1().Services(namespace).Create(context.TODO(), &source, v1.CreateOptions{})
				if err != nil {
					log.Printf("创建目标 service 失败 %s: %v", source.Name, err)
					continue
				}
				log.Printf("在目标集群中成功创建了 service %s", source.Name)
			} else {
				log.Printf("获取目标 Service 失败 %s: %v", source.Name, err)
			}
			continue
		}
		// 处理 svc
		targetRes, err := targetCli.CoreV1().Services(namespace).Update(context.TODO(), targetService, v1.UpdateOptions{})
		if err != nil {
			log.Printf("更新目标 service 失败 %s: %v", targetRes.Name, err)
		} else {
			log.Printf("成功同步 service %s", targetRes.Name)
		}

	}

	return nil
}
