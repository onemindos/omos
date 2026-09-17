# Machine Learning Curriculum

## Foundation & Roadmap

This curriculum is anchored to practical, project-based learning. It bridges theory and implementation through hands-on work with real datasets.

---

## Phase 1: Fundamentals (Weeks 1-4)

### 1.1 Mathematics Foundation
- **Linear Algebra**: Vectors, matrices, operations (NumPy)
- **Statistics**: Distributions, hypothesis testing, probability
- **Calculus**: Derivatives, gradients (backpropagation intuition)

**Resources:**
- 3Blue1Brown: Essence of Linear Algebra
- StatQuest with Josh Starmer: Statistics & ML fundamentals
- Khan Academy: Calculus

**Output:** Math notebook with worked examples

---

### 1.2 Python Essentials
- Data types, control flow, functions
- **NumPy**: Arrays, vectorization
- **Pandas**: DataFrames, data manipulation
- **Matplotlib / Seaborn**: Visualization

**Resources:**
- Python docs
- NumPy/Pandas official tutorials
- Kaggle notebooks (beginner tier)

**Output:** Python scripts for data wrangling

---

## Phase 2: ML Fundamentals (Weeks 5-8)

### 2.1 Supervised Learning
- **Regression**: Linear, polynomial, Ridge/Lasso
- **Classification**: Logistic regression, decision trees, random forests, SVM
- Train/test split, cross-validation, hyperparameter tuning

**Resources:**
- Scikit-learn documentation
- Andrew Ng's ML course (free on Coursera)
- "Hands-On Machine Learning" by Aurélien Géron

**Output:** Classification & regression models on toy datasets

---

### 2.2 Unsupervised Learning
- **Clustering**: K-means, hierarchical, DBSCAN
- **Dimensionality Reduction**: PCA, t-SNE
- Evaluation metrics (silhouette score, elbow method)

**Resources:**
- Scikit-learn clustering docs
- Kaggle: Unsupervised Learning challenges

**Output:** Clustering analysis notebook

---

## Phase 3: Practical Projects (Weeks 9-12)

### **Project 1: House Price Prediction**
- Dataset: Boston Housing or Kaggle House Prices
- Task: Regression (predict prices)
- Skills: EDA, feature engineering, linear/tree models, evaluation
- Deliverable: Jupyter notebook + trained model

### **Project 2: Iris/Titanic Classification**
- Dataset: UCI Iris or Kaggle Titanic
- Task: Binary/multi-class classification
- Skills: Data cleaning, class imbalance, logistic regression, ensemble methods
- Deliverable: Model pipeline + accuracy metrics

### **Project 3: Customer Segmentation**
- Dataset: Kaggle Mall Customers or similar
- Task: Clustering
- Skills: Preprocessing, K-means, PCA, visualization
- Deliverable: Segment profiles + business insights

### **Project 4: Sentiment Analysis**
- Dataset: Movie reviews or Twitter data
- Task: Text classification
- Skills: Tokenization, TF-IDF or embeddings, naive Bayes/SVM
- Deliverable: Sentiment classifier

### **Project 5: Time Series Forecasting**
- Dataset: Stock prices, weather, or sales data
- Task: Predict future values
- Skills: ARIMA/SARIMA, trend/seasonality decomposition, evaluation (MAE, RMSE)
- Deliverable: Forecasting model + confidence intervals

---

## Phase 4: Advanced Topics (Weeks 13-16)

### 4.1 Deep Learning Basics
- **Neural Networks**: Perceptron, MLP, backpropagation
- **Frameworks**: TensorFlow / PyTorch basics
- **CNNs**: For image classification (MNIST, CIFAR-10)
- **RNNs/LSTMs**: For sequence problems

**Resources:**
- TensorFlow/PyTorch official tutorials
- Fast.ai: Practical Deep Learning for Coders
- Jeremy Howard's videos

**Output:** CNN classifier, LSTM time series model

---

### 4.2 Production & MLOps
- Model serialization (pickle, ONNX)
- API deployment (Flask, FastAPI)
- Docker containerization
- Model versioning & monitoring

**Resources:**
- Made With ML: MLOps fundamentals
- Full Stack Python: Web framework patterns

**Output:** REST API serving a trained model

---

## Phase 5: Specialization (Ongoing)

Choose one or more paths:

### **Computer Vision**
- Image classification, object detection, segmentation
- Tools: OpenCV, Detectron2, YOLOv5
- Project: Build a custom object detector

### **NLP**
- Word embeddings, transformers, LLMs
- Tools: NLTK, spaCy, Hugging Face
- Project: Fine-tune a transformer for text classification

### **Reinforcement Learning**
- MDPs, policy gradients, Q-learning
- Tools: Gym, Stable Baselines3
- Project: Train an agent on a game/simulation

### **Time Series & Forecasting**
- Advanced methods: Prophet, AutoARIMA, neural forecasting
- Project: Multi-step forecasting pipeline

---

## Learning Stack & Tools

### Essential
- **Python 3.10+**
- **Jupyter Notebook / VS Code**
- **Scikit-learn** (classical ML)
- **NumPy / Pandas** (data wrangling)
- **Matplotlib / Seaborn** (visualization)

### Intermediate
- **TensorFlow / PyTorch** (deep learning)
- **XGBoost / LightGBM** (gradient boosting)
- **Optuna / Hyperopt** (hyperparameter tuning)

### Advanced
- **Hugging Face Transformers** (NLP)
- **FastAPI / Flask** (deployment)
- **Docker** (containerization)
- **MLflow / Weights & Biases** (experiment tracking)

---

## Datasets & Communities

### Free Datasets
- **Kaggle**: 10,000+ datasets + competitions
- **UCI Machine Learning Repository**: Classic datasets
- **Google Dataset Search**: Federated search
- **OpenML**: Curated ML datasets
- **Papers with Code**: Datasets linked to research

### Communities & Support
- **Kaggle**: Competitions, notebooks, forums
- **Fast.ai**: Forum, courses, community
- **Reddit**: r/MachineLearning, r/learnmachinelearning
- **GitHub**: Open-source ML projects

---

## Assessment & Milestones

| Milestone | Criteria | Timeline |
|-----------|----------|----------|
| **Phase 1 Complete** | Math & Python fundamentals mastered | Week 4 |
| **Phase 2 Complete** | Built 2 supervised + 1 unsupervised model | Week 8 |
| **5 Projects Done** | All projects completed with docs & results | Week 12 |
| **Portfolio Ready** | GitHub repo with 5 polished projects | Week 16 |
| **Specialization** | Deep dive into chosen area (CV/NLP/RL) | Ongoing |

---

## Resource Index

### Anchor Resource
- **Instagram Post**: @chrispathway - "5 Machine Learning Projects You Can Build This Weekend"
  - Link: https://www.instagram.com/p/DaVkcTIiO2y/
  - Offers: Code + datasets for 5 beginner ML projects
  - Call-to-action: Comment "Project" to receive build guide

### Key References
1. **Andrew Ng's Machine Learning Specialization** (Coursera) - Foundations
2. **"Hands-On Machine Learning"** (Géron) - Theory + practice
3. **Fast.ai** - Practical deep learning from practitioners
4. **Kaggle Learn** - Free micro-courses + datasets
5. **Papers with Code** - Research + implementations

---

## Next Steps

1. **Week 1**: Set up Python environment (venv/Anaconda), install core libraries
2. **Week 2-4**: Complete math & Python modules
3. **Week 5-8**: Work through ML fundamentals with Scikit-learn
4. **Week 9-12**: Build the 5 projects from the @chrispathway resource (use the datasets/code they provide)
5. **Week 13+**: Specialize based on interest (DL, NLP, CV, RL)

---

**Created:** 2026-07-09  
**Status:** Active  
**Last Updated:** 2026-07-09
