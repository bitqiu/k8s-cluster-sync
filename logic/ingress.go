package logic

import (
	"context"
	"fmt"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func SyncIngress(sourceCli, targetCli *kubernetes.Clientset, namespace string) error {
	sourceList, err := sourceCli.NetworkingV1().Ingresses(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list source ingress: %v", err)
	}

	for _, source := range sourceList.Items {
		target, err := targetCli.NetworkingV1().Ingresses(namespace).Get(context.TODO(), source.Name, v1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				targetCreateRes, err := targetCli.NetworkingV1().Ingresses(namespace).Create(context.TODO(), target, v1.CreateOptions{})
				if err != nil {
					log.Printf("创建目标 ingress 失败 %s: %v", targetCreateRes.Name, err)
					continue
				}
				log.Printf("在目标集群中成功创建了 ingress %s", targetCreateRes.Name)
			} else {
				log.Printf("获取目标 ingress 失败 %s: %v", source.Name, err)
			}
			continue
		}

		// 如果存在，则更新目标集群中的 Deployment
		targetUpdateRes, err := targetCli.NetworkingV1().Ingresses(namespace).Update(context.TODO(), &source, v1.UpdateOptions{})
		if err != nil {
			log.Printf("更新目标 ingress 失败 %s: %v", targetUpdateRes.Name, err)
		} else {
			log.Printf("成功同步 ingress %s", targetUpdateRes.Name)
		}
	}

	return nil
}
