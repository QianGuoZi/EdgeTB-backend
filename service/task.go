package service

import (
	"EdgeTB-backend/dal"
	"log"
)

type GetAllTaskResponse struct {
	ID          int64  `json:"id"`
	DatasetName string `json:"datasetName"`
	ConfigId    int64  `json:"configId"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

func GetAllTask(username, projectName string) ([]GetAllTaskResponse, error) {
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[GetAllTask] 服务获取用户id失败")
		return nil, err
	}
	// 获取项目id
	projectId, err := dal.GetProjectId(projectName, userId)
	if err != nil {
		log.Printf("[GetAllTask] 服务获取项目id失败")
		return nil, err
	}

	// 获取任务列表
	tasks, err := dal.GetAllTask(projectId)
	if err != nil {
		log.Printf("[GetAllTask] 服务获取任务列表失败")
		return nil, err
	}

	var res []GetAllTaskResponse
	for _, task := range tasks {
		dataset, err := dal.GetDatasetDetail(int(task.DatasetId))
		if err != nil {
			log.Printf("[GetAllTask] 服务获取数据集名称失败")
			return nil, err
		}
		res = append(res, GetAllTaskResponse{
			ID:          task.Id,
			DatasetName: dataset.DatasetName,
			ConfigId:    task.ConfigId,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

type AddTaskRequest struct {
	DatasetId       int64         `json:"datasetId"`
	ConfigId        int64         `json:"configId"`
	DatasetSplitter *UploadedFile `json:"datasetSplitter"`
}

func AddTask(username, projectName string, req AddTaskRequest) (int64, error) {
	userId, err := dal.GetUserId(username)
	if err != nil {
		log.Printf("[AddTask] 服务获取用户id失败")
		return 0, err
	}
	// 获取项目id
	projectId, err := dal.GetProjectId(projectName, userId)
	if err != nil {
		log.Printf("[AddTask] 服务获取项目id失败")
		return 0, err
	}

	splitterFile := dal.File{
		Name: req.DatasetSplitter.FileName,
		Url:  req.DatasetSplitter.URL,
		Size: req.DatasetSplitter.Size,
	}
	fileId, err := dal.AddFile(splitterFile)
	if err != nil {
		log.Printf("[AddTask] 服务创建文件失败")
		return 0, err
	}

	// 创建任务
	taskId, err := dal.AddTask(dal.Task{
		Status:                dal.TaskStatusCreated,
		DatasetId:             req.DatasetId,
		DatasetSplitterFileId: fileId,
		ConfigId:              req.ConfigId,
		ProjectId:             projectId,
	})
	if err != nil {
		log.Printf("[AddTask] 服务创建任务失败")
		return 0, err
	}

	return taskId, nil
}
