# CES-27 - LAB 1

Instruções para configuração do ambiente: [Setup](SETUP.md)  
Para saber mais sobre a utilização do Git, acessar: [Git Passo-a-Passo](GIT.md)  
Para saber mais sobre a entrega do trabalho, acessar: [Entrega](ENTREGA.md)  

## Informações Gerais:

O bloco a seguir indica que um comando deve ser executado na linha de comando:
```shell
folder$ git --version
```
> git version 2.9.2.windows.1

**folder** é a pasta onde o comando deve ser executado. Caso esteja vazio, o comando deve ser executado na pasta raiz do projeto (ces27-lab1).  
**$** indica o começo do comando a ser executado.  
**git version 2.9.2.windows.1** é a saída no console.  

## MapReduce

Nesse laboratório será implementada uma versão simplificada do modelo de programação MapReduce proposto por Jeffrey Dean e Sanjay Ghemawat no paper [MapReduce: Simplified Data Processing on Large Clusters](http://research.google.com/archive/mapreduce.html).
Implementações de MapReduce normalmente rodam em grandes clusters executando paralelamente e de forma distribuída milhares de tarefas. Além disso, é fundamental que as frameworks sejam capazes de lidar com problemas básicos presentes em sistemas distribuídos: falhas e capacidade de operar em larga escala.

O MapReduce é inspirado pelas funções map e reduce comumente utilizadas em programação funcional.

**Map**  
A função de Map é responsável por fazer um mapeamento dos dados de entrada em uma estrutura do tipo lista de chaves/valores. Esse mapeamento é executado de forma paralela para diversos dados de entradas gerando uma lista por *job* executado.

**Reduce**  
A função reduce recebe os resultados das funções de Map e faz uma sumarização dos valores obtidos, gerando uma nova lista de chaves/valores.

Além dessas, outras funções são normalmente utilizadas: Shuffle(ou Partition), que vai distribuir os resultados dos maps em várias partições (múltiplos reduces), garantindo que os dados que similares serão alocados na mesma partição para ser executado por um mesmo Reduce.

## Parte II - Execução Distribuída com Tolerância a Falhas

Nesta segunda parte da prática, temos pronto um processo de MapReduce que já opera de forma distribuída.

Da prática anterior, devemos ter prontas as funções de Map, Reduce e Split. Copiar as duas primeiras no arquivo **wordcount/wordcount.go** e a última no arquivo **wordcount/data.go**.

### Executando o Código
Para executar o código, precisamos executar multiplos processos. O primeiro deles são os clientes. Note que o valor de -port deve ser alterado para cada cliente. Neste lab usaremos a port 5000 para o Master e 5001+ para os Workers.
O segundo é o Master.

Executando Worker (aumentar -port para o worker seguinte):
```bash
wordcount$ go run main.go wordcount.go data.go -mode distributed -type worker -port 5001
```
> 2016/09/01 00:52:58 Running in distributed mode.  
> 2016/09/01 00:52:58 NodeType: worker  
> 2016/09/01 00:52:58 Address: localhost  
> 2016/09/01 00:52:58 Port: 5001  
> 2016/09/01 00:52:58 Master: localhost:5000  
> 2016/09/01 00:52:58 Running Worker on localhost:5001  
> 2016/09/01 00:52:58 Registering with Master  

Executando o Master:
```bash
wordcount$ go run main.go wordcount.go data.go -mode distributed -type master
```
> 2016/09/01 00:55:38 Running in distributed mode.  
> 2016/09/01 00:55:38 NodeType: master  
> 2016/09/01 00:55:38 Reduce Jobs: 5  
> 2016/09/01 00:55:38 Address: localhost  
> 2016/09/01 00:55:38 Port: 5000  
> 2016/09/01 00:55:38 File: files/pg1342.txt  
> 2016/09/01 00:55:38 Chunk Size: 102400  
> 2016/09/01 00:55:38 Running Master on localhost:5000  
> 2016/09/01 00:55:38 Scheduling Worker.RunMap operations  
> 2016/09/01 00:55:38 Accepting connections on 127.0.0.1:5000  


Quando os dois processos tiverem sido inicializados, temos as seguintes saídas:

Master:
> 2016/09/01 00:57:35 Registering worker '0' with hostname 'localhost:5001'  
> 2016/09/01 00:57:35 Running Worker.RunMap (ID: '0' File: 'map\map-0' Worker: '0')  
> 2016/09/01 00:57:36 Running Worker.RunMap (ID: '1' File: 'map\map-1' Worker: '0')  
> 2016/09/01 00:57:37 Running Worker.RunMap (ID: '2' File: 'map\map-2' Worker: '0')  
> 2016/09/01 00:57:38 Running Worker.RunMap (ID: '3' File: 'map\map-3' Worker: '0')  
> 2016/09/01 00:57:39 Running Worker.RunMap (ID: '4' File: 'map\map-4' Worker: '0')  
> 2016/09/01 00:57:40 Running Worker.RunMap (ID: '5' File: 'map\map-5' Worker: '0')  
> 2016/09/01 00:57:41 Running Worker.RunMap (ID: '6' File: 'map\map-6' Worker: '0')  
> 2016/09/01 00:57:42 Running Worker.RunMap (ID: '7' File: 'map\map-7' Worker: '0')  
> 2016/09/01 00:57:43 Worker.RunMap operations completed  
> 2016/09/01 00:57:43 Scheduling Worker.RunReduce operations  
> 2016/09/01 00:57:43 Running Worker.RunReduce (ID: '0' File: 'reduce\reduce-0' Worker: '0')  
> 2016/09/01 00:57:44 Running Worker.RunReduce (ID: '1' File: 'reduce\reduce-1' Worker: '0')  
> 2016/09/01 00:57:45 Running Worker.RunReduce (ID: '2' File: 'reduce\reduce-2' Worker: '0')  
> 2016/09/01 00:57:46 Running Worker.RunReduce (ID: '3' File: 'reduce\reduce-3' Worker: '0')  
> 2016/09/01 00:57:47 Running Worker.RunReduce (ID: '4' File: 'reduce\reduce-4' Worker: '0')  
> 2016/09/01 00:57:48 Worker.RunReduce operations completed  
> 2016/09/01 00:57:48 Closing Remote Workers.  
> 2016/09/01 00:57:49 Done.  

Worker:
> 2016/09/01 00:57:35 Registered. WorkerId: 0  
> 2016/09/01 00:57:35 Accepting connections on 127.0.0.1:5001  
> 2016/09/01 00:57:36 Running map id: 0, path: map\map-0  
> 2016/09/01 00:57:37 Running map id: 1, path: map\map-1  
> 2016/09/01 00:57:38 Running map id: 2, path: map\map-2  
> 2016/09/01 00:57:39 Running map id: 3, path: map\map-3  
> 2016/09/01 00:57:40 Running map id: 4, path: map\map-4  
> 2016/09/01 00:57:41 Running map id: 5, path: map\map-5  
> 2016/09/01 00:57:42 Running map id: 6, path: map\map-6  
> 2016/09/01 00:57:43 Running map id: 7, path: map\map-7  
> 2016/09/01 00:57:44 Running reduce id: 0, path: reduce\reduce-0  
> 2016/09/01 00:57:45 Running reduce id: 1, path: reduce\reduce-1  
> 2016/09/01 00:57:46 Running reduce id: 2, path: reduce\reduce-2  
> 2016/09/01 00:57:47 Running reduce id: 3, path: reduce\reduce-3  
> 2016/09/01 00:57:48 Running reduce id: 4, path: reduce\reduce-4  
> 2016/09/01 00:57:49 Done.  

### Atividade
Executar com 1 master e 1 worker.
Executar com 1 master e 2 workers. Para isso você deve inicializar primeiramente os Workers e só então inicializar o Master.

### Discussão
Nessa execução, a tarefa de MapReduce deve estar sendo concluída corretamente. É importante observar que a execução não é determinística.

O grande problema desse modelo é que ele não é resiliente a falhas. Quando executado em larga escala, milhares de nós são alocados e executam o processo de Worker. Com muitos nós alocados, é comum que ocorram falhas e que alguns desses processos parem.

### Introduzindo Falhas
Vamos introduzir agora um processo worker no qual vamos induzir uma falha.

Para isso, inicialize o worker usando o seguinte comando:
```bash
wordcount$ go run main.go wordcount.go data.go -mode distributed -type worker -port 5001 -fail 3
```
> 2016/09/01 01:09:51 Running in distributed mode.  
> 2016/09/01 01:09:51 NodeType: worker  
> 2016/09/01 01:09:51 Address: localhost  
> 2016/09/01 01:09:51 Port: 5001  
> 2016/09/01 01:09:51 Master: localhost:5000  
> 2016/09/01 01:09:51 Induced failure  
> 2016/09/01 01:09:51 After 3 operations  
> 2016/09/01 01:09:51 Running Worker on localhost:5001  
> 2016/09/01 01:09:51 Registering with Master  

Nesse caso, esse worker vai falhar na execução da terceira tarefa.

Executando o Master junto com esse Worker defeituoso, teremos o seguinte resultado:
```go
wordcount$ go run main.go data.go wordcount.go -mode distributed -type master
```
> 2016/09/01 01:10:34 Running in distributed mode.  
> 2016/09/01 01:10:34 NodeType: master  
> 2016/09/01 01:10:34 Reduce Jobs: 5  
> 2016/09/01 01:10:34 Address: localhost  
> 2016/09/01 01:10:34 Port: 5000  
> 2016/09/01 01:10:34 File: files/pg1342.txt  
> 2016/09/01 01:10:34 Chunk Size: 102400  
> 2016/09/01 01:10:34 Running Master on localhost:5000  
> 2016/09/01 01:10:34 Scheduling Worker.RunMap operations  
> 2016/09/01 01:10:34 Accepting connections on 127.0.0.1:5000  
> 2016/09/01 01:10:36 Registering worker '0' with hostname 'localhost:5001'  
> 2016/09/01 01:10:36 Running Worker.RunMap (ID: '0' File: 'map\map-0' Worker: '0')  
> 2016/09/01 01:10:37 Running Worker.RunMap (ID: '1' File: 'map\map-1' Worker: '0')  
> 2016/09/01 01:10:38 Operation Worker.RunMap '1' Failed. Error: read tcp 127.0.0.1:58204->127.0.0.1:5001: wsarecv: An existing connection was forcibly closed by the remote host.  

Note que a execução vai ficar travada (deadlock). Para entender melhor vamos olhar o seguintes trechos de código:


```go
// master.go
type Master struct {
    (...)
    // Workers handling
    workersMutex sync.Mutex
    workers      map[int]*RemoteWorker
    totalWorkers int // Used to generate unique ids for new workers

    idleWorkerChan   chan *RemoteWorker
    failedWorkerChan chan *RemoteWorker
    (...)
}
```

Todas essas estruturas são utilizadas para gerenciar os Workers no Master. 

A primeira coisa que acontece quando um worker é inicializado é a chamada da função Register no Master:

```go
// master_rpc.go
func (master *Master) Register(args *RegisterArgs, reply *RegisterReply) error {
    (...)
    master.workersMutex.Lock()

    newWorker = &RemoteWorker{master.totalWorkers, args.WorkerHostname, WORKER_IDLE}
    master.workers[newWorker.id] = newWorker
    master.totalWorkers++

    master.workersMutex.Unlock()

    master.idleWorkerChan <- newWorker
    (...)
}
```

Observe que quando um novo Worker se registra no Master, ele adiciona esse worker na lista **master.workers**. Para evitar problemas de sincronia no acesso dessa estrutura usamos o mutex **master.workersMutex**.

O ponto mais interessante entretanto, é a última linha:
```go
master.idleWorkerChan <- newWorker
```

Nessa linha, o worker que acabou de se registrado é colocado no canal de Workers disponíveis. A saída desse canal está em:

```go
// master_scheduler.go
func (master *Master) schedule(task *Task, proc string, filePathChan chan string) {
    (...)
    counter = 0
    for filePath = range filePathChan {
        operation = &Operation{proc, counter, filePath}
        counter++

        worker = <-master.idleWorkerChan
        wg.Add(1)
        go master.runOperation(worker, operation, &wg)
    }

    wg.Wait()
    (...)
}
```

Essa é a função que aloca as operações nos workers.

Mas porque essas coisas estão acontecendo em paralelo?

```go
// mapreduce.go
    (...)
    go master.acceptMultipleConnections()
    go master.handleFailingWorkers()

    // Schedule map operations
    master.schedule(task, "Worker.RunMap", task.InputFilePathChan)
    (...)
```

Observe que na linha master.acceptMultipleConnections() é colocado em uma Goroutine separada. Isso é o equivalente a colocar a execução desse método numa Thread separada.
Já na linha master.schedule(...), fazemos a execuçao na goroutine atual.

Dessa forma essas operações estão acontecendo concorrentemente, e o master.idleWorkerChan é o canal de comunicação entre elas. Quando a operação Register escreve no canal, a operação schedule é informado de que um novo Worker está disponível e continua a execução.

Por último, para entender o Deadlock, vamos olhar o código a seguir:
```go
func (master *Master) runOperation(remoteWorker *RemoteWorker, operation *Operation, wg *sync.WaitGroup) {
    (...)
    err = remoteWorker.callRemoteWorker(operation.proc, args, new(struct{}))

    if err != nil {
        log.Printf("Operation %v '%v' Failed. Error: %v\n", operation.proc, operation.id, err)
        wg.Done()
        master.failedWorkerChan <- remoteWorker
    } else {
        wg.Done()
        master.idleWorkerChan <- remoteWorker
    }
    (...)
}
```

Quando um worker completa uma operação corretamente, ele cai no else acima e o worker que a executou é colocado no canal master.idleWorkerChan (e que vai ser lido pelo scheduler mostrado anteriormente). Entretanto, no caso de falha, o worker é colocado no canal master.failedWorkerChan
### Atividade ###