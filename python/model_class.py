import pandas as pd
from sklearn.model_selection import train_test_split, GridSearchCV
from sklearn.pipeline import Pipeline
from sklearn.compose import ColumnTransformer
from sklearn.impute import SimpleImputer
from sklearn.preprocessing import OneHotEncoder
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import accuracy_score, classification_report
import joblib
import sys


def load_data(file_path):
    if file_path.endswith(".csv"):
        return pd.read_csv(file_path)
    elif file_path.endswith(".xls") or file_path.endswith(".xlsx"):
        return pd.read_excel(file_path)
    elif file_path.endswith(".pkl"):
        return pd.read_pickle(file_path)
    else:
        raise ValueError(
            "Unsupported file type. Supported types are csv, xls, xlsx, and pkl."
        )


def train_model(data, target_column, model_path, algorithm="random_forest"):
    X = data.drop(columns=[target_column])
    y = data[target_column]

    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    categorical_features = X.select_dtypes(include=["object"]).columns
    numeric_features = X.select_dtypes(include=["int", "float"]).columns

    categorical_transformer = Pipeline(
        steps=[
            ("imputer", SimpleImputer(strategy="most_frequent")),
            ("onehot", OneHotEncoder(handle_unknown="ignore")),
        ]
    )

    numeric_transformer = Pipeline(
        steps=[
            ("imputer", SimpleImputer(strategy="mean")),
        ]
    )

    preprocessor = ColumnTransformer(
        transformers=[
            ("num", numeric_transformer, numeric_features),
            ("cat", categorical_transformer, categorical_features),
        ]
    )

    if algorithm == "random_forest":
        classifier = RandomForestClassifier()
        param_grid = {
            # "classifier__n_estimators": [100, 200],
            # "classifier__min_samples_split": [5, 10],
        }
    elif algorithm == "logistic_regression":
        classifier = LogisticRegression()
        param_grid = (
            {}
        )  # "classifier__C": [0.1, 1, 10], "classifier__max_iter": [100, 300]}
    else:
        raise ValueError(
            "Unsupported algorithm. Choose 'random_forest' or 'logistic_regression'."
        )

    clf = Pipeline(steps=[("preprocessor", preprocessor), ("classifier", classifier)])

    grid_search = GridSearchCV(clf, param_grid, cv=5, scoring="accuracy")
    grid_search.fit(X_train, y_train)

    y_pred = grid_search.predict(X_test)
    accuracy = accuracy_score(y_test, y_pred)
    report = classification_report(y_test, y_pred, output_dict=True)

    precision = report["weighted avg"]["precision"]
    recall = report["weighted avg"]["recall"]
    f1_score = report["weighted avg"]["f1-score"]

    print("Accuracy:", accuracy)
    print(f"Precision: {precision:.2f}")
    print(f"Recall: {recall:.2f}")
    print(f"F1-score: {f1_score:.2f}")

    joblib.dump(grid_search.best_estimator_, model_path)

    return grid_search.best_estimator_


def main():
    if len(sys.argv) != 5:
        print(
            "Usage: python model_class.py <target_column> <input_file_path> <output_file_path>"
        )
        sys.exit(1)

    algorithm = sys.argv[1]
    target_column = sys.argv[2]
    input_file_path = sys.argv[3]
    output_file_path = sys.argv[4]

    data = load_data(input_file_path)
    train_model(data, target_column, algorithm=algorithm, model_path=output_file_path)


if __name__ == "__main__":
    main()
